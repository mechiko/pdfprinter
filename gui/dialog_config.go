package gui

import (
	"pdfprinter/gui/dconfig"
	"pdfprinter/reductor"
)

func (a *GuiApp) onConfig() {
	model, _ := GetModel()
	if model == nil {
		return
	}
	data := dconfig.ConfigDialogData{
		PerLabel:     model.PerLabel,
		MarkTemplate: model.MarkTemplate,
		ChunkSize:    model.ChunkSize,
	}
	dlg := dconfig.NewConfigDialog(&data)
	dlg.ShowModal()
	if data.Ok {
		// model.PerLabel = data.PerLabel
		// model.MarkTemplate = data.MarkTemplate
		model.ChunkSize = data.ChunkSize
		if err := model.SyncToStore(a); err != nil {
			a.Logger().Errorf("диалог onConfig синхронизация модели %v", err)
		}
		err := reductor.Instance().SetModel(model, false)
		if err != nil {
			a.Logger().Errorf("dialog onConfig set reductor error %v", err)
		}
	}
}
