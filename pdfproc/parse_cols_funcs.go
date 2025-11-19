package pdfproc

import (
	"fmt"
	"pdfprinter/domain"
	"strconv"
	"strings"

	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/mechiko/utility"
)

func (p *pdfProc) addDatamaxColumn(col core.Col, colTempl *domain.RowPrimitive, ciss []*utility.CisInfo) error {
	var cis *utility.CisInfo
	indexCis, err := strconv.ParseInt(colTempl.DataMatrix, 10, 64)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if int(indexCis) > len(ciss)-1 {
		return fmt.Errorf("datamatrix index %s out of bounds", colTempl.DataMatrix)
	}
	cis = ciss[indexCis]
	fnc := cis.FNC1()
	img, err := dmImg(fnc)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	col.Add(
		image.NewFromBytes(img, extension.Png, colTempl.PropsRect()),
	)
	return nil
}

func (p *pdfProc) addBarColumn(col core.Col, colTempl *domain.RowPrimitive, ciss []*utility.CisInfo) error {
	cis := ciss[0]
	switch colTempl.Bar {
	case "ean13":
		ean13 := strings.Trim(cis.Gtin, "0")
		col.Add(
			code.NewBar(ean13, colTempl.PropsBar()),
		)
	}
	return nil
}

func (p *pdfProc) addJpgColumn(col core.Col, colTempl *domain.RowPrimitive, ciss []*utility.CisInfo) error {
	if p.assets != nil {
		img, err := p.assets.Jpg(colTempl.Image)
		if err != nil {
			return fmt.Errorf("page image assets %w", err)
		}
		if len(img) == 0 {
			return fmt.Errorf("page image assets empty for %q", colTempl.Image)
		}
		col.Add(
			image.NewFromBytes(img, colTempl.ImageExt, colTempl.PropsRect()),
		)
	} else {
		return fmt.Errorf("page image assets not available (assets is nil) for %q", colTempl.Image)
	}
	return nil
}

func (p *pdfProc) addStringColumn(col core.Col, colTempl *domain.RowPrimitive, ciss []*utility.CisInfo) error {
	cis := ciss[0]
	party, err := p.vars.Get("party")
	if err != nil {
		return fmt.Errorf("page vars party get error %w", err)
	}
	idx, err := p.vars.Get("idx")
	if err != nil {
		return fmt.Errorf("page vars idx get error %w", err)
	}
	value := strings.ReplaceAll(colTempl.Value, "@party", party)
	value = strings.ReplaceAll(value, "@idx", idx)
	ean13 := strings.Trim(cis.Gtin, "0")
	value = strings.ReplaceAll(value, "@ean", ean13)
	col.Add(text.New(value, colTempl.PropsText()))
	return nil
}

func (p *pdfProc) addArrayStringColumn(col core.Col, colTempl *domain.RowPrimitive, ciss []*utility.CisInfo) error {
	comps := make([]core.Component, 0)
	cis := ciss[0]
	party, err := p.vars.Get("party")
	if err != nil {
		return fmt.Errorf("page vars party get error %w", err)
	}
	idx, err := p.vars.Get("idx")
	if err != nil {
		return fmt.Errorf("page vars idx get error %w", err)
	}
	for _, val := range colTempl.Values {
		value := ""
		if val.Value != "" {
			value = strings.ReplaceAll(val.Value, "@party", party)
			value = strings.ReplaceAll(value, "@idx", idx)
			ean13 := strings.Trim(cis.Gtin, "0")
			ean13 = fmt.Sprintf("%s  %s  %s", ean13[:1], ean13[1:7], ean13[7:])
			value = strings.ReplaceAll(value, "@ean", ean13)
			comps = append(comps, text.New(value, val.PropsText()))
		}
		if val.DataMatrix != "" {
			fnc := cis.FNC1()
			img, err := dmImg(fnc)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			comps = append(comps, image.NewFromBytes(img, extension.Png, colTempl.PropsRect()))
		}
		if val.Bar != "" {
			switch val.Bar {
			case "ean13":
				ean13 := strings.Trim(cis.Gtin, "0")
				comps = append(comps, code.NewBar(ean13, colTempl.PropsBar()))
			case "ean13b":
				ean13 := strings.Trim(cis.Gtin, "0")
				img, err := barImg(ean13)
				if err != nil {
					return fmt.Errorf("ean13 bar %w", err)
				}
				comps = append(comps, image.NewFromBytes(img, extension.Jpg, val.PropsRect()))
			case "ean13svg":
				ean13 := strings.Trim(cis.Gtin, "0")
				img, err := svgImg(ean13)
				if err != nil {
					return fmt.Errorf("ean13 svgImg %w", err)
				}
				comps = append(comps, image.NewFromBytes(img, extension.Jpg, val.PropsRect()))
			case "ean13j":
				img, err := p.assets.Jpg("gtin2")
				if err != nil {
					return fmt.Errorf("page image assets %w", err)
				}
				if len(img) == 0 {
					return fmt.Errorf("page image assets empty for %q", colTempl.Image)
				}
				comps = append(comps, image.NewFromBytes(img, extension.Jpg, val.PropsRect()))
			case "ean13p":
				img, err := p.assets.Png(cis.Gtin)
				if err != nil {
					return fmt.Errorf("page image assets %w", err)
				}
				if len(img) == 0 {
					return fmt.Errorf("page image assets empty for %q", colTempl.Image)
				}
				comps = append(comps, image.NewFromBytes(img, extension.Png, val.PropsRect()))
			}
		}
	}
	col.Add(comps...)
	return nil
}
