package pdfproc

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/ean"
)

// генератор EAN
func barImg(code string) ([]byte, error) {
	bcImg, err := ean.Encode(code)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	scaled, err := barcode.Scale(bcImg, 100, 30)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var b bytes.Buffer
	jpeg.Encode(&b, scaled, nil)
	png.Encode(&b, scaled)
	return b.Bytes(), nil
}

// генератор PNG
func (p *pdfProc) dmImg(code string) ([]byte, error) {
	if code == "" {
		return nil, fmt.Errorf("generate datamatrix error: code is empty")
	}
	bcImg, err := datamatrix.Encode(code)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	dx := bcImg.Bounds().Dx()
	dy := bcImg.Bounds().Dy()
	scaled, err := barcode.Scale(bcImg, dx, dy)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	scaledCheck, err := barcode.Scale(bcImg, 50, 50)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	err = p.reader(code, scaledCheck)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var b bytes.Buffer
	// jpeg.Encode(&b, scaled, nil)
	png.Encode(&b, scaled)
	return forceTo8BitPNG(b.Bytes()), nil
}

func forceTo8BitPNG(data []byte) []byte {
	src, _ := png.Decode(bytes.NewReader(data))
	b := src.Bounds()
	dst := image.NewNRGBA(b)
	draw.Draw(dst, b, src, b.Min, draw.Src)
	buf := new(bytes.Buffer)
	png.Encode(buf, dst)
	return buf.Bytes()
}

// генератор PNG
func dmImgJpg(code string, size int) ([]byte, error) {
	// fmt.Printf("KM len %d\n", len(code))
	bcImg, err := datamatrix.Encode(code)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	scaled, err := barcode.Scale(bcImg, size, size)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var b bytes.Buffer
	jpeg.Encode(&b, scaled, nil)
	return b.Bytes(), nil
}
