package pdfkm

import (
	"fmt"
	"log"
	"path/filepath"
	"pdfprinter/domain/models/application"
	"pdfprinter/pdfproc"
	"slices"

	"github.com/mechiko/utility"
)

func (k *Pdf) DocumentChunk(model *application.Application, chProgress chan float64, step float64, chunk *ChunkPack, file string) error {
	if k.templateDatamatrix == nil {
		return fmt.Errorf("Error pdfkm datamatrix template is nil ")
	}
	pdfDocument, err := pdfproc.New(k, k.assets)
	if err != nil {
		return fmt.Errorf("Error create pdfproc: %v", err)
	}
	if err := pdfDocument.BuildMaroto(k.templateDatamatrix); err != nil {
		return fmt.Errorf("build maroto error %w", err)
	}
	cises := chunk.Cis
	party := model.Party
	if r := []rune(party); len(r) > 2 {
		party = string(r[:2])
	}
	if model.PerLabel <= 0 {
		return fmt.Errorf("PerLabel must be > 0, got %d", model.PerLabel)
	}
	switch {
	case model.PerLabel == 1:
		for _, ciss := range cises {
			// генерируем этикетку по одной
			k.iChunkAll++
			k.iChunkCis++
			// генерируем КМ
			err = pdfDocument.SetVars("party", party)
			if err != nil {
				log.Fatalf("vars add party %v", err)
			}
			err = pdfDocument.SetVars("idx", fmt.Sprintf("%06d", k.iChunkCis))
			if err != nil {
				log.Fatalf("vars add idx %v", err)
			}
			if err := pdfDocument.AddPageByTemplate(k.templateDatamatrix, []*utility.CisInfo{ciss}); err != nil {
				return fmt.Errorf("add datamatrix KM in page (idx %d): %w", k.iChunkAll, err)
			}

			k.SendProgress(chProgress, step*float64(k.iChunkAll))
		}
	case model.PerLabel > 1:
		iLabel := 0
		packs := slices.Chunk(cises, model.PerLabel)
		for ciss := range packs {
			// генерируем этикетку по несколько штук
			k.iChunkAll += model.PerLabel
			k.iChunkCis += model.PerLabel
			iLabel++
			// генерируем КМ
			err = pdfDocument.SetVars("party", party)
			if err != nil {
				log.Fatalf("vars add party %v", err)
			}
			err = pdfDocument.SetVars("idx", fmt.Sprintf("%06d", k.iChunkCis))
			if err != nil {
				log.Fatalf("vars add idx %v", err)
			}
			if err := pdfDocument.AddPageByTemplate(k.templateDatamatrix, ciss); err != nil {
				return fmt.Errorf("add datamatrix KM in page %d: %w", iLabel, err)
			}

			k.SendProgress(chProgress, step*float64(k.iChunkAll))
		}
	}

	err = pdfDocument.DocumentGenerate()
	if err != nil {
		return fmt.Errorf("генерация пдф блока ошибка %w", err)
	}
	err = pdfDocument.PdfDocumentSave(filepath.Join(model.FileBasePath, file))
	if err != nil {
		return fmt.Errorf("генерация пдф блока save error %q: %w", file, err)
	}
	return nil
}
