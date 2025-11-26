package pdfproc

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strconv"

	"github.com/mechiko/dmxing"
	"github.com/mechiko/dmxing/datamatrix"
)

func (p *pdfProc) reader(code string, png bytes.Buffer) error {
	img, _, err := image.Decode(&png)
	if err != nil {
		return fmt.Errorf("Unable to decode to DataMatrix code to bitmap: %s", err)
	}
	// prepare BinaryBitmap
	bmp, err := dmxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return fmt.Errorf("Unable to decode to DataMatrix code to bitmap: %s", err)
	}

	// decode image
	datamatrixReader := datamatrix.NewDataMatrixReader()
	result, err := datamatrixReader.Decode(bmp, nil)
	if err != nil {
		return fmt.Errorf("Unable to decode to DataMatrix code to bitmap: %s", err)
	}
	hx := strconv.Quote(result.GetText())
	codeQuoted := strconv.Quote(code)
	if hx == codeQuoted {
		return nil
	}
	return fmt.Errorf("code [%s] not equal png decode [%s]", codeQuoted, hx)

}
