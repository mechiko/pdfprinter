package pdfproc

import (
	"bytes"
	"fmt"
	imgstd "image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"pdfprinter/domain"
	"slices"
	"strings"

	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/page"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/juliankoehn/barcode"
	svgbar "github.com/juliankoehn/barcode"
	"github.com/juliankoehn/barcode/ean"
	"github.com/mechiko/utility"
)

func (p *pdfProc) Page(t *domain.MarkTemplate, ciss []*utility.CisInfo, party string, idx string) (core.Page, error) {
	if len(ciss) != 4 {
		return nil, fmt.Errorf("len cis must be 4")
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
			// одна строка автороу
			row1 := rowTempl[0]
			if row1.Value == "" {
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
					pg.Add(
						row.New(row1.RowHeight).Add(
							text.NewCol(12, row1.Value, row1.PropsText()),
						),
					)
				}
			}
		case len(rowTempl) > 1:
			cols := make([]core.Col, len(rowTempl))
			// строки с колонками
			for i, rowSingle := range rowTempl {
				cols[i] = col.New(rowSingle.ColWidth)
				if rowSingle.DataMatrix != "" {
					var cis *utility.CisInfo
					switch rowSingle.DataMatrix {
					case "0":
						cis = ciss[0]
					case "1":
						cis = ciss[1]
					case "2":
						cis = ciss[2]
					case "3":
						cis = ciss[3]
					default:
						cis = ciss[0]
					}
					fnc := cis.FNC1()
					cols[i].Add(
						code.NewMatrix(fnc, rowSingle.PropsRect()),
					)
				}
				if rowSingle.Bar != "" {
					cis := ciss[0]
					switch rowSingle.Bar {
					case "ean13":
						ean13 := strings.Trim(cis.Gtin, "0")
						cols[i].Add(
							code.NewBar(ean13, rowSingle.PropsBar()),
						)
					}
				}
				if rowSingle.Image != "" {
					if p.assets != nil {
						img, err := p.assets.Jpg(rowSingle.Image)
						if err != nil {
							return nil, fmt.Errorf("page image assets error %w", err)
						}
						if len(img) == 0 {
							return nil, fmt.Errorf("page image assets empty for %q", rowSingle.Image)
						}
						cols[i].Add(
							image.NewFromBytes(img, rowSingle.ImageExt, rowSingle.PropsRect()),
						)
					} else {
						return nil, fmt.Errorf("page image assets not available (assets is nil) for %q", rowSingle.Image)
					}
				}
				if rowSingle.Value != "" {
					cis := ciss[0]
					value := strings.ReplaceAll(rowSingle.Value, "@party", party)
					value = strings.ReplaceAll(value, "@idx", idx)
					ean13 := strings.Trim(cis.Gtin, "0")
					value = strings.ReplaceAll(value, "@ean", ean13)
					cols[i].Add(text.New(value, rowSingle.PropsText()))
				} else {
					if len(rowSingle.Values) > 0 {
						comps := make([]core.Component, 0)
						for _, val := range rowSingle.Values {
							cis := ciss[0]
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
								comps = append(comps, code.NewMatrix(cis.FNC1(), val.PropsRect()))
							}
							if val.Bar != "" {
								switch val.Bar {
								case "ean13":
									ean13 := strings.Trim(cis.Gtin, "0")
									comps = append(comps, code.NewBar(ean13, rowSingle.PropsBar()))
								case "ean13b":
									ean13 := strings.Trim(cis.Gtin, "0")
									img, err := barImg(ean13)
									if err != nil {
										return pg, fmt.Errorf("ean13 bar error %w", err)
									}
									comps = append(comps, image.NewFromBytes(img, extension.Jpg, val.PropsRect()))
								case "ean13svg":
									ean13 := strings.Trim(cis.Gtin, "0")
									img, err := svgImg(ean13)
									if err != nil {
										return pg, fmt.Errorf("ean13 svgImg error %w", err)
									}
									comps = append(comps, image.NewFromBytes(img, extension.Jpg, val.PropsRect()))
								case "ean13j":
									img, err := p.assets.Jpg("gtin2")
									if err != nil {
										return nil, fmt.Errorf("page image assets error %w", err)
									}
									if len(img) == 0 {
										return nil, fmt.Errorf("page image assets empty for %q", rowSingle.Image)
									}
									comps = append(comps, image.NewFromBytes(img, extension.Jpg, val.PropsRect()))
								case "ean13p":
									img, err := p.assets.Png(cis.Gtin)
									if err != nil {
										return nil, fmt.Errorf("page image assets error %w", err)
									}
									if len(img) == 0 {
										return nil, fmt.Errorf("page image assets empty for %q", rowSingle.Image)
									}
									comps = append(comps, image.NewFromBytes(img, extension.Png, val.PropsRect()))
								}
							}
						}
						cols[i].Add(comps...)
					}
				}
			}
			pg.Add(
				row.New(rowTempl[0].RowHeight).Add(cols...),
			)
		default:
		}
	}
	return pg, nil
}

func barImg(code string) ([]byte, error) {
	bcImg, err := ean.Encode(code)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	scaled, err := barcode.Scale(bcImg, 100, 25)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var b bytes.Buffer
	jpeg.Encode(&b, scaled, nil)
	png.Encode(&b, scaled)
	return b.Bytes(), nil
}

func svgImg(code string) ([]byte, error) {
	var b bytes.Buffer
	width := 2
	height := 35
	color := "black"
	showCode := false
	which := "EAN13"
	extension := ".jpg"
	if extension == ".png" || extension == ".jpg" || extension == ".jpeg" || extension == ".gif" {
		var img *imgstd.RGBA
		if extension == ".png" {
			_, img = svgbar.GetBarcodeFile(code, which, width, height, color, showCode, false, true)
		} else {
			_, img = svgbar.GetBarcodeFile(code, which, width, height, color, showCode, false, false)
		}
		switch extension {
		case ".png":
			png.Encode(&b, img)
		case ".jpg", ".jpeg":
			jpeg.Encode(&b, img, &jpeg.Options{Quality: 100})
		case ".gif":
			gif.Encode(&b, img, nil)
		}
	}
	return b.Bytes(), nil
}
