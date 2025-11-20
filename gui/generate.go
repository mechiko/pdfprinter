package gui

import (
	"fmt"
	"path/filepath"
	"pdfprinter/pdfkm"
	"pdfprinter/reductor"
	"time"

	"github.com/mechiko/utility"
)

// кнопка Пуск
// запускать в отдельном поток от tk9
func (a *GuiApp) generate() {
	defer func() {
		a.SendIsProcess(false)
	}()

	a.SendIsProcess(true)
	model, err := SyncModel(a)
	if err != nil {
		a.logerr("gui generate sync model: ", err)
		return
	}
	a.SendLogClear()

	tMark := fmt.Sprintf("выбран шаблон печати КМ: %s", model.MarkTemplate)
	a.SendLog(tMark)

	// проверяем файлы
	cisCode, sizeFile, err := pdfkm.CheckFile(model.FileCIS)
	if err != nil {
		a.logerr("ошибка проверки файлов: ", err)
		return
	}
	sizeCode, err := sizeCis(cisCode)
	if err != nil {
		a.logerr("ошибка проверки размера кода CIS : ", err)
		return
	}

	tSize := fmt.Sprintf("загружено КМ %d, размер DM %d ", sizeFile, sizeCode)
	a.SendLog(tSize)

	a.SendLog("обрабатываем файл...")
	pdfGenerator, err := pdfkm.New(a)
	if err != nil {
		a.logerr("генерация пдф:", err)
		return
	}

	if err := pdfGenerator.ReadCIS(model); err != nil {
		model.FileCIS = ""
		a.logerr("ошибка загрузки файла:", err)
		return
	}
	// запрашиваем имя выходного файла и пути
	filenameCis := filepath.Base(model.FileCIS)
	extensionCis := filepath.Ext(model.FileCIS)
	fileBaseName := filenameCis[:len(filenameCis)-len(extensionCis)]
	fileNamePdf := utility.TimeFileName(fileBaseName) + ".pdf"
	fileNamePdfSelect, err := utility.DialogSaveFile(utility.Pdf, fileNamePdf, ".")
	if err != nil {
		a.logerr("генерация пдф: выбор пути для сохранения PDF", err)
	} else if fileNamePdfSelect != "" {
		fileNamePdf = fileNamePdfSelect
	}
	model.SetFileBase(fileNamePdf)
	if err := reductor.Instance().SetModel(model, false); err != nil {
		a.Logger().Errorf("генерация пдф: ошибка сохранения модели %s", err.Error())
		a.SendError(fmt.Sprintf("генерация пдф: ошибка сохранения модели %s", err.Error()))
		return
	}

	startgen := time.Now()
	// сплит на блоки по chunksize
	err = pdfGenerator.ChunkSplit(model)
	if err != nil {
		a.logerr("генерация пдф: сплит на блоки", err)
		return
	}

	// здесь генерируем документ ПДФ целиком
	err = pdfGenerator.Document(model, a.progresCh)
	if err != nil {
		a.logerr("генерация пдф: ", err)
		return
	}

	since := fmt.Sprintf("%v", time.Since(startgen))
	a.SendLog("сгенерированы файлы: за " + since)
	for _, file := range pdfGenerator.Files() {
		a.SendLog(file)
	}
	utility.MessageBox("Сообщение", "формирование пдф завершено успешно")
	a.SendFinish()
}
