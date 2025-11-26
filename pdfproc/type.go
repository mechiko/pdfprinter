package pdfproc

import (
	"fmt"
	"pdfprinter/assets"
	"pdfprinter/domain"

	"github.com/johnfercher/maroto/v2/pkg/core"
)

type pdfProc struct {
	domain.Apper
	maroto   core.Maroto
	assets   *assets.Assets
	document core.Document
	debug    bool
	height   float64
	width    float64
	vars     *Vars
}

func New(app domain.Apper, assets *assets.Assets) (*pdfProc, error) {
	if app == nil {
		return nil, fmt.Errorf("app is nil")
	}
	if assets == nil {
		return nil, fmt.Errorf("assets is nil")
	}
	p := &pdfProc{
		Apper:  app,
		assets: assets,
		vars:   NewVars(),
	}
	return p, nil
}

func (p *pdfProc) SetVars(name, value string) error {
	return p.vars.Add(name, value)
}
