package gui

import (
	_ "embed"
	"fmt"
	"pdfprinter/domain"
	"sync"
	"time"

	tk "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

const (
	tick = 10 * time.Millisecond
)

type LogMsg struct {
	Error bool
	Msg   string
}

//go:embed 192.png
var ico []byte

type GuiApp struct {
	domain.Apper
	icon *tk.Img

	buttonFrame *tk.TFrameWidget
	inputFrame  *tk.TFrameWidget
	logFrame    *tk.TFrameWidget

	startButton  *tk.TButtonWidget
	exitButton   *tk.TButtonWidget
	configButton *tk.TButtonWidget
	logCh        chan LogMsg
	// stateFinishOpenXlsx   chan struct{}
	stateFinish           chan struct{}
	stateStart            chan struct{}
	logClear              chan struct{}
	stateSelectedCisFile  chan string
	stateSelectedKiguFile chan string
	stateIsProcess        chan bool
	stateFinishDebug      chan struct{}
	yscroll               *tk.Window
	logText               *tk.TextWidget

	// processing *processing.Processing
	// pdf     *pdfkm.Pdf
	fileLblCis *tk.TLabelWidget
	fileBtnCis *tk.TButtonWidget

	progres   *tk.TProgressbarWidget
	progresCh chan float64
	isProcess bool
	chunkSize *tk.TEntryWidget

	datamatrixCombo *tk.TComboboxWidget

	lock sync.Mutex
}

func New(app domain.Apper) (*GuiApp, error) {
	a := &GuiApp{
		Apper: app,
		// pdf:   p,
	}
	a.logCh = make(chan LogMsg, 100)
	a.stateFinish = make(chan struct{}, 1)
	a.stateFinishDebug = make(chan struct{}, 1)
	a.stateStart = make(chan struct{}, 1)
	a.icon = tk.NewPhoto(tk.Data(ico))
	a.progresCh = make(chan float64, 1)
	a.logClear = make(chan struct{}, 1)
	a.stateSelectedCisFile = make(chan string, 1)
	a.stateSelectedKiguFile = make(chan string, 1)
	a.stateIsProcess = make(chan bool, 1)

	tk.App.IconPhoto(a.icon)
	tk.ErrorMode = tk.CollectErrors
	tk.App.WmTitle("Формирование ПДФ КМ")
	tk.WmProtocol(tk.App, "WM_DELETE_WINDOW", a.onQuitApp)
	if err := tk.ActivateTheme("azure light"); err != nil {
		a.Logger().Errorf("gui theme %s", err.Error())
	}
	tk.InitializeExtension("autoscroll")

	model, err := GetModel()
	if err != nil {
		return nil, fmt.Errorf("gui new get model %w", err)
	}
	a.makeWidgets(model)
	a.makeLayout()
	a.makeBindings()
	// start ticker only after widgets/layout are ready
	tk.NewTicker(tick, a.tick)
	return a, nil
}

func (a *GuiApp) Run() {
	tk.App.Center()
	// before Run() is called and before the GUI event loop starts.
	// If openFileCis accesses GUI widgets or sends to channels
	// that are processed by the ticker, this could lead to race conditions or deadlocks
	model, _ := GetModel()
	if model.FileCIS != "" {
		go a.openFileCis(model.FileCIS)
	}

	tk.WmDeiconify(tk.App)
	tk.App.Wait()
}

func (a *GuiApp) logg(s, e string) {
	blue := "color1"
	red := "color2"
	if s != "" {
		s += "\n"
	}
	if e != "" {
		e += "\n"
	}
	a.logText.Configure(tk.State("normal"))
	a.logText.Insert(tk.END, s, blue, e, red)
	a.logText.See("end")
	a.logText.Configure(tk.State("disabled"))
}

func (a *GuiApp) onQuitApp() {
	// If this field is modified by other goroutines (e.g., in the ticker or worker threads),
	// this creates a data race
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.isProcess {
		a.logg("", "выход из программы ограничен, запущена обработка")
		return
	}
	tk.Destroy(tk.App)
}
