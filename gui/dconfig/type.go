package dconfig

import (
	"strconv"

	tk "modernc.org/tk9.0"
)

type ConfigDialogData struct {
	Ok           bool
	PerLabel     int
	MarkTemplate string
	ChunkSize    int
}

type ConfigDialog struct {
	data *ConfigDialogData
	win  *tk.ToplevelWidget

	perLabel *tk.TEntryWidget

	buttonFrame  *tk.TFrameWidget
	inputFrame   *tk.TFrameWidget
	okButton     *tk.TButtonWidget
	cancelButton *tk.TButtonWidget

	datamatrixCombo *tk.TComboboxWidget
	chunkSize       *tk.TEntryWidget
}

func NewConfigDialog(data *ConfigDialogData) *ConfigDialog {
	dlg := &ConfigDialog{data: data}
	dlg.win = tk.App.Toplevel()
	dlg.win.WmTitle("Config")
	// tk.WmAttributes(dlg.win, tk.Type("dialog")) // TODO
	tk.WmProtocol(dlg.win.Window, tk.WM_DELETE_WINDOW, dlg.onCancel)

	dlg.makeWidgets()
	dlg.makeLayout()
	dlg.makeBindings()
	return dlg
}

func (me *ConfigDialog) onOk() {
	me.data.Ok = true
	if per, err := strconv.ParseInt(me.perLabel.Textvariable(), 10, 64); err == nil {
		me.data.PerLabel = int(per)
	}
	tk.Destroy(me.win)
}

func (me *ConfigDialog) onCancel() {
	tk.Destroy(me.win)
}

func (me *ConfigDialog) ShowModal() {
	me.win.Raise(tk.App)
	tk.Focus(me.win)
	// tk.Focus(me.percentSpinbox)
	tk.GrabSet(me.win)
	me.win.Center().Wait()
}
