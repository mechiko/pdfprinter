package gui

import (
	"fmt"
	"pdfprinter/pdfkm"
	"pdfprinter/reductor"

	"github.com/mechiko/utility"
)

// кнопка Пуск
// запускать в отдельном поток от tk9
func (a *GuiApp) generate() {
	defer func() {
		a.stateIsProcess <- false
	}()
	a.stateIsProcess <- true

	model, err := GetModel()
	if err != nil {
		a.Logger().Errorf("gui generate %s", err.Error())
		a.SendError(fmt.Sprintf("gui generate %s", err.Error()))
		return
	}
	if err := model.SyncToStore(a); err != nil {
		a.Logger().Errorf("ошибка синхронизации модели в настройки программы %s", err.Error())
		a.SendError(fmt.Sprintf("ошибка синхронизации модели в настройки программы %s", err.Error()))
		return
	}
	// сохраняем модель по ошибке
	logerr := func(s string, err error) {
		if err := reductor.Instance().SetModel(model, false); err != nil {
			a.Logger().Errorf("gui generate setmodel %s", err.Error())
			a.SendError(fmt.Sprintf("gui generate setmodel  %s", err.Error()))
			return
		}
		if err != nil {
			a.Logger().Errorf("%s %s", s, err.Error())
			a.SendError(fmt.Sprintf("%s %s", s, err.Error()))
			a.stateStart <- struct{}{}
		}
	}
	a.logClear <- struct{}{}
	tMark := fmt.Sprintf("выбран шаблон печати КМ: %s", model.MarkTemplate)
	a.SendLog(tMark)

	// проверяем файлы
	sizeCis, err := pdfkm.CheckFile(model.FileCIS)
	if err != nil {
		logerr("ошибка проверки файлов: ", err)
		return
	}
	tSize := fmt.Sprintf("загружено КМ: %d", sizeCis)
	a.SendLog(tSize)

	a.SendLog("обрабатываем файл...")
	pdfGenerator, err := pdfkm.New(a)
	if err != nil {
		logerr("генерация пдф:", err)
		return
	}
	a.SendLog("считываем файл КМ")
	if err := pdfGenerator.ReadCIS(model); err != nil {
		model.FileCIS = ""
		logerr("ошибка загрузки файла:", err)
		return
	}
	// запрашиваем имя выходного файла и пути
	fileNamePdf := utility.TimeFileName(model.MarkTemplate) + ".pdf"
	fileNamePdfSelect, err := utility.DialogSaveFile(utility.Pdf, fileNamePdf, ".")
	if err != nil {
		logerr("генерация пдф: выбор пути для сохранения PDF", err)
	} else if fileNamePdfSelect != "" {
		fileNamePdf = fileNamePdfSelect
	}
	model.SetFileBase(fileNamePdf)
	if err := reductor.Instance().SetModel(model, false); err != nil {
		a.Logger().Errorf("генерация пдф: ошибка сохранения модели %s", err.Error())
		a.SendError(fmt.Sprintf("генерация пдф: ошибка сохранения модели %s", err.Error()))
		return
	}
	// сплит на блоки по chunksize
	err = pdfGenerator.ChunkSplit(model)
	if err != nil {
		logerr("генерация пдф: сплит на блоки", err)
		if model != nil && model.FileCIS != "" {
			a.stateSelectedCisFile <- model.FileCIS
		}
		return
	}

	// здесь генерируем документ ПДФ целиком
	err = pdfGenerator.Document(model, a.progresCh)
	if err != nil {
		logerr("генерация пдф: документ ошибка", err)
		if model != nil && model.FileCIS != "" {
			a.stateSelectedCisFile <- model.FileCIS
		}
		return
	}

	a.SendLog("сгенерированы файлы:")
	for _, file := range pdfGenerator.Files() {
		a.SendLog(file)
	}
	utility.MessageBox("Сообщение", "формирование пдф завершено успешно")
	a.stateFinish <- struct{}{}
}
