package pdfproc

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strconv"

	"github.com/mechiko/dmxing"
	"github.com/mechiko/dmxing/datamatrix"
)

func (p *pdfProc) reader(code string, img image.Image) error {
	// prepare BinaryBitmap
	bmp, err := dmxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return fmt.Errorf("Unable to convert to DataMatrix-code to bitmap: %s", err)
	}

	// decode image
	datamatrixReader := datamatrix.NewDataMatrixReader()
	result, err := datamatrixReader.Decode(bmp, nil)
	if err != nil {
		return fmt.Errorf("Unable to convert to DataMatrix-code to bitmap: %s", err)
	}
	hx := strconv.Quote(result.GetText())
	cisQuted := strconv.Quote(code)
	if hx == cisQuted {
		return nil
	}
	return fmt.Errorf("code [%s] not equal png decode [%s]", cisQuted, hx)

}
