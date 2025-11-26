package pdfproc

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strconv"

	"github.com/mechiko/dmxing"
	"github.com/mechiko/dmxing/datamatrix"
)

func (p *pdfProc) reader(code string, img image.Image) bool {
	// prepare BinaryBitmap
	bmp, err := dmxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		p.Logger().Errorf("Unable to convert to DataMatrix-code to bitmap: %s", err)
		return false
	}

	// decode image
	datamatrixReader := datamatrix.NewDataMatrixReader()
	result, err := datamatrixReader.Decode(bmp, nil)
	if err != nil {
		p.Logger().Errorf("Unable to convert to DataMatrix-code to bitmap: %s", err)
		return false
	}
	hx := strconv.Quote(result.GetText())
	cisQuted := strconv.Quote(code)
	if hx == cisQuted {
		return true
	}
	return false
}
