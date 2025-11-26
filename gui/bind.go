package gui

import (
	"fmt"
	"pdfprinter/domain"
	"pdfprinter/domain/models/application"
	"pdfprinter/reductor"
	"strconv"

	tk "modernc.org/tk9.0"
)

func (a *GuiApp) makeBindings() {
	// tk.Bind(tk.App, "<Escape>", tk.Command(a.onQuitApp))
	// tk.Bind(tk.App, "<<ComboboxSelected>>", tk.Command(func() {
	// 	model, err := GetModel()
	// 	if err != nil {
	// 		a.Logger().Errorf("gui new get model %w", err)
	// 		return
	// 	}
	// 	model.Magazin = a.magazinCombo.Textvariable()
	// 	// a.magazinCombo.Configure(tk.Textvariable(model.Magazin))
	// 	if _, ok := model.Reestr[model.Magazin]; ok {
	// 		a.SendLog(fmt.Sprintf("выбран магазин %s", model.Magazin))
	// 		a.startButton.Configure(tk.State("enabled"))
	// 		reductor.Instance().SetModel(model, false)
	// 	} else {
	// 		a.startButton.Configure(tk.State("disabled"))
	// 	}
	// }))
	tk.Bind(a.datamatrixCombo, "<<ComboboxSelected>>", tk.Command(func() {
		appModel, err := GetModel()
		if err != nil {
			a.logerr("get model", err)
			return
		}
		tmpl, err := appModel.SetMarkTemplate(a.datamatrixCombo.Textvariable())
		if err != nil {
			a.logerr("model set marktemplate", err)
			return
		}
		str := fmt.Sprintf("шаблон: размер H:%0.0f W:%0.0f марок на этикетке:%d", tmpl.PageHeight, tmpl.PageWidth, tmpl.KmPlace)
		a.SendLog(str)
		err = reductor.Instance().SetModel(appModel, false)
		if err != nil {
			a.logerr("update model", err)
			return
		}
	}))
	tk.Bind(a.chunkSize, "<KeyRelease>", tk.Command(func(e *tk.Event) {
		txt := a.chunkSize.Textvariable()
		n, err := strconv.Atoi(txt)
		if err != nil || n <= 0 {
			a.Logger().Warnf("invalid ChunkSize: %q", txt)
			a.chunkSize.Configure(tk.Textvariable(""))
			return
		}
		mdl, err := reductor.Instance().Model(domain.Application)
		if err != nil {
			a.Logger().Errorf("get model error: %v", err)
			return
		}
		appModel, ok := mdl.(*application.Application)
		if !ok {
			a.Logger().Errorf("bad type model application")
			return
		}
		appModel.ChunkSize = n
		reductor.Instance().SetModel(appModel, false)
	}))
}
