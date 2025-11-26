package pdfproc

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"

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
	delta := (dx * 4) + (dx / 3)
	scaledCheck, err := barcode.Scale(bcImg, delta, delta)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var b, bc bytes.Buffer
	err = png.Encode(&b, scaled)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	err = png.Encode(&bc, scaledCheck)
	if err != nil {
		return nil, fmt.Errorf("encode 8-bit png: %w", err)
	}
	if p.DebugMode() {
		if err := os.WriteFile("file.png", b.Bytes(), 0666); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		if err := os.WriteFile("check.png", bc.Bytes(), 0666); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}
	err = p.reader(code, bc)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	out, err := forceTo8BitPNG(b.Bytes())
	if err != nil {
		return nil, fmt.Errorf("normalize datamatrix png: %w", err)
	}
	return out, nil
}

func forceTo8BitPNG(data []byte) ([]byte, error) {
	src, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode png: %w", err)
	}
	b := src.Bounds()
	dst := image.NewNRGBA(b)
	draw.Draw(dst, b, src, b.Min, draw.Src)
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, dst); err != nil {
		return nil, fmt.Errorf("encode 8-bit png: %w", err)
	}
	return buf.Bytes(), nil
}

// генератор PNG
func dmImgJpg(code string, size int) ([]byte, error) {
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
