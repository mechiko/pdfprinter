package gui

import (
	"fmt"
	"path/filepath"
	"pdfprinter/pdfkm"
	"pdfprinter/reductor"

	"github.com/boombuler/barcode/datamatrix"
)

// должна выполнятся как gorutine
func (a *GuiApp) openFileCis(file string) {

	model, err := GetModel()
	if err != nil {
		a.logerr("gui openFile", err)
		return
	}
	model.FileCIS = file
	model.SetFileBase(file)
	err = reductor.Instance().SetModel(model, false)
	if err != nil {
		a.logerr("ошибка записи модели в редуктор:", err)
		return
	}
	// a.SendLogClear()
	if file != "" {
		base := filepath.Base(file)
		if len(base) > 50 {
			base = fmt.Sprintf("%.30s...%s", base, base[len(base)-10:])
		}
		a.SendLog("проверяем файл КМ " + base)
		cis, size, err := pdfkm.CheckFile(file)
		if err != nil {
			a.logerr("ошибка проверки файла: ", err)
			return
		}
		a.SendLog(fmt.Sprintf("считано %d КМ", size))
		sizeCis, err := sizeCis(cis)
		if err != nil {
			a.logerr("ошибка проверки размера кода CIS : ", err)
			return
		}
		a.SendLog(fmt.Sprintf("размер DM %d ед", sizeCis))
	}
	// устанавливаем состояни отображения gui
	a.SendSelectedCisFile(model.FileCIS)
}

// генератор PNG
func sizeCis(code string) (int, error) {
	if code == "" {
		return 0, fmt.Errorf("code is empty")
	}
	bcImg, err := datamatrix.Encode(code)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	dx := bcImg.Bounds().Dx()
	return dx, nil
}
