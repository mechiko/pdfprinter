package gui

import "fmt"

// вызывать из gorutine
// из основного потока вызывать только как go
func (a *GuiApp) SendError(s string) {
	msg := LogMsg{
		Error: true,
		Msg:   s,
	}
	a.Logger().Error(s)
	select {
	case a.logCh <- msg:
		// message sent
	default:
		// message dropped, log this event
		a.Logger().Warn("Failed to send error message to GUI logCh: channel full %s", s)
	}
}

// вызывать из gorutine
// из основного потока вызывать только как go
func (a *GuiApp) SendLog(s string) {
	msg := LogMsg{
		Error: false,
		Msg:   s,
	}
	select {
	case a.logCh <- msg:
		// message sent
	default:
		// message dropped
		a.Logger().Debug("Failed to send log message to GUI logCh: channel full %s", s)
	}
}

func (a *GuiApp) SendProgress(f float64) {
	select {
	case a.progresCh <- f:
		// message sent
	default:
		// message dropped
		a.Logger().Debug("Failed to send to GUI progresCh: channel full caller:%s", callerFunctionName())
	}
}

func (a *GuiApp) SendIsProcess(f bool) {
	select {
	case a.stateIsProcess <- f:
		// message sent
	default:
		// message dropped
		a.Logger().Debug("Failed to send to GUI stateIsProcess: channel full caller:%s", callerFunctionName())
	}
}

func (a *GuiApp) SendStart() {
	select {
	case a.stateStart <- struct{}{}:
		// message sent
	default:
		// message dropped
		a.Logger().Debug("Failed to send to GUI stateStart: channel full caller:%s", callerFunctionName())
	}
}

func (a *GuiApp) SendLogClear() {
	select {
	case a.logClear <- struct{}{}:
		// message sent
	default:
		// message dropped
		a.Logger().Debug("Failed to send to GUI logClear: channel full caller:%s", callerFunctionName())
	}
}

func (a *GuiApp) logerr(s string, err error) {
	if err != nil {
		a.Logger().Errorf("%s %s", s, err.Error())
		a.SendError(fmt.Sprintf("%s %s", s, err.Error()))
		a.SendStart()
	}
}

func (a *GuiApp) SendFinish() {
	select {
	case a.stateFinish <- struct{}{}:
		// message sent
	default:
		// message dropped
		a.Logger().Debug("Failed to send to GUI stateFinish: channel full caller:%s", callerFunctionName())
	}
}

func (a *GuiApp) SendSelectedCisFile(f string) {
	select {
	case a.stateSelectedCisFile <- f:
		// message sent
	default:
		// message dropped
		a.Logger().Debug("Failed to send to GUI stateSelectedCisFile: channel full caller:%s", callerFunctionName())
	}
}
