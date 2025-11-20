package pdfproc

import (
	"fmt"
	"pdfprinter/domain"
	"slices"

	"github.com/johnfercher/maroto/v2/pkg/components/page"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/mechiko/utility"
)

// в станицу передаются только cis для нее
func (p *pdfProc) Page(t *domain.MarkTemplate, ciss []*utility.CisInfo) (core.Page, error) {
	if len(ciss) != t.KmPlace {
		return nil, fmt.Errorf("мест в шаблоне %d а марок передано %d", t.KmPlace, len(ciss))
	}
	pg := page.New()
	rowKeys := make([]string, 0, len(t.Rows))
	for k := range t.Rows {
		rowKeys = append(rowKeys, k)
	}
	slices.Sort(rowKeys)
	for _, rowKey := range rowKeys {
		rowTempl := t.Rows[rowKey]
		switch {
		case len(rowTempl) == 0:
		case len(rowTempl) == 1:
			row1 := rowTempl[0]
			// одна строка автороу
			if err := p.parseSingleRow(pg, row1, ciss); err != nil {
				return nil, fmt.Errorf("parse single row error %w", err)
			}
		case len(rowTempl) > 1:
			if err := p.parseColsRow(pg, rowTempl, ciss); err != nil {
				return nil, fmt.Errorf("parse cols row error %w", err)
			}
		default:
		}
	}
	return pg, nil
}
