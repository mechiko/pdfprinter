package gui

import (
	"fmt"
	"path/filepath"

	tk "modernc.org/tk9.0"
)

func (a *GuiApp) tick() {
	select {
	case s := <-a.logCh:
		if s.Error {
			a.logg("", s.Msg)
		} else {
			a.logg(s.Msg, "")
		}
	case <-a.logClear:
		a.logText.Configure(tk.State("normal"))
		a.logText.Delete("1.0", tk.END)
		a.logText.Configure(tk.State("disabled"))
	case v := <-a.progresCh:
		a.progres.Configure(tk.Value(v))
	case <-a.stateStart:
		// состояние начала
		a.progres.Configure(tk.Value(0))
	case <-a.stateFinish:
		// создание удаляет предыдущий канал и будет собран мусорщиком
		// a.progresCh = make(chan float64)
		a.progres.Configure(tk.Value(0))
	case <-a.stateFinishDebug:
	case file := <-a.stateSelectedCisFile:
		label := ""
		if file != "" {
			base := filepath.Base(file)
			if len(base) > 50 {
				label = fmt.Sprintf("%.30s...%s", base, base[len(base)-10:])
			} else {
				label = base
			}
		}
		a.fileLblCis.Configure(tk.Txt(label))
	case a.isProcess = <-a.stateIsProcess:
		if a.isProcess {
			a.fileBtnCis.Configure(tk.State("disabled"))
			a.startButton.Configure(tk.State("disabled"))
			a.exitButton.Configure(tk.State("disabled"))
		} else {
			a.fileBtnCis.Configure(tk.State("enabled"))
			a.startButton.Configure(tk.State("enabled"))
			a.exitButton.Configure(tk.State("enabled"))
		}
	default:
	}
}
