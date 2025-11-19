package pdfproc

import (
	"fmt"
	"pdfprinter/domain"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/mechiko/utility"
)

// если не указывается высота то вставляется авто роу не важно на остальное
// и текст из value
func (p *pdfProc) parseSingleRow(pg core.Page, row1 *domain.RowPrimitive, ciss []*utility.CisInfo) error {
	if row1.Value == "" && row1.DataMatrix == "" {
		// пустая строка с высотой
		pg.Add(
			row.New(row1.RowHeight).Add(),
		)
	} else {
		if row1.RowHeight == 0 {
			pg.Add(
				text.NewAutoRow(row1.Value, row1.PropsText()),
			)
		} else {
			cis := ciss[0]
			colNew := col.New(12)
			if row1.DataMatrix != "" {
				fnc := cis.FNC1()
				img, err := dmImg(fnc)
				if err != nil {
					return fmt.Errorf("%w", err)
				}
				colNew.Add(image.NewFromBytes(img, extension.Png, row1.PropsRect()))
				pg.Add(
					row.New(row1.RowHeight).Add(
						colNew,
					),
				)

			} else {
				pg.Add(
					row.New(row1.RowHeight).Add(
						text.NewCol(12, row1.Value, row1.PropsText()),
					),
				)
			}
		}
	}
	return nil
}
