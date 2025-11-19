package pdfproc

import (
	"fmt"
	"pdfprinter/domain"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/mechiko/utility"
)

// ciss КМ передаются в строку колонок или одна или несколько
func (p *pdfProc) parseColsRow(pg core.Page, colsTempl []*domain.RowPrimitive, ciss []*utility.CisInfo) error {
	cols := make([]core.Col, 0)
	// строки с колонками
	for i, colSingle := range colsTempl {
		if colSingle.ColWidth == 0 {
			continue
		}
		colNew := col.New(colSingle.ColWidth)
		cols = append(cols, colNew)

		switch {
		case colSingle.DataMatrix != "":
			if err := p.addDatamaxColumn(colNew, colSingle, ciss); err != nil {
				return fmt.Errorf("generate col %d %w", i, err)
			}
		case colSingle.Bar != "":
			if err := p.addBarColumn(colNew, colSingle, ciss); err != nil {
				return fmt.Errorf("generate col %d %w", i, err)
			}
		case colSingle.Image != "":
			if err := p.addJpgColumn(colNew, colSingle, ciss); err != nil {
				return fmt.Errorf("generate col %d %w", i, err)
			}
		case colSingle.Value != "":
			if err := p.addStringColumn(colNew, colSingle, ciss); err != nil {
				return fmt.Errorf("generate col %d %w", i, err)
			}
		case len(colSingle.Values) > 0:
			if err := p.addArrayStringColumn(colNew, colSingle, ciss); err != nil {
				return fmt.Errorf("generate col %d %w", i, err)
			}
		}
	}
	pg.Add(
		row.New(colsTempl[0].RowHeight).Add(cols...),
	)
	return nil
}
