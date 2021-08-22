package flusher

import (
	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"github.com/ozonva/ova-algorithm-api/internal/repo"
)

// Flusher - интерфейс для сброса задач в хранилище
type Flusher interface {
	Flush(algorithms []algorithm.Algorithm) []algorithm.Algorithm
}

// NewFlusher возвращает Flusher с поддержкой батчевого сохранения
func NewFlusher(chunkSize int, algorithmRepo repo.Repo) Flusher {
	return &flusher{
		chunkSize:     chunkSize,
		algorithmRepo: algorithmRepo,
	}
}

type flusher struct {
	chunkSize     int
	algorithmRepo repo.Repo
}

func (f *flusher) Flush(algorithmsInput []algorithm.Algorithm) []algorithm.Algorithm {
	bulks := algorithm.SplitAlgorithmsToBulks(algorithmsInput, uint(f.chunkSize))

	var failedAlgorithms []algorithm.Algorithm

	for i := 0; i < len(bulks); i++ {
		if err := f.algorithmRepo.AddAlgorithms(bulks[i]); err != nil {
			failedAlgorithms = append(failedAlgorithms, bulks[i]...)
		}
	}

	return failedAlgorithms
}
