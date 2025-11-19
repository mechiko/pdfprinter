package pdfkm

import (
	"fmt"
	"pdfprinter/domain/models/application"
	"slices"
)

func (k *Pdf) ChunkSplit(model *application.Application) error {
	if model.ChunkSize <= 0 || model.PerLabel <= 0 {
		return fmt.Errorf("некорректные параметры разбиения: ChunkSize=%d, PerPallet=%d", model.ChunkSize, model.PerLabel)
	}
	countCisChunk := model.ChunkSize * model.PerLabel
	countCIS := 0
	chunksCIS := slices.Chunk(k.Cis, countCisChunk)
	k.OrderChunks = make([]string, 0)
	for chunk := range chunksCIS {
		chunkPack := &ChunkPack{
			Cis: chunk,
		}
		fileChunk := fmt.Sprintf("%06d-%06d_%s.pdf", countCIS*countCisChunk+1, ((countCIS + 1) * countCisChunk), model.FileBaseName)
		k.OrderChunks = append(k.OrderChunks, fileChunk)
		k.Chunks[fileChunk] = chunkPack
		countCIS++
	}
	return nil
}
