package flusher

import (
	"fmt"
	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"github.com/ozonva/ova-algorithm-api/internal/repo"
)

// Flusher - интерфейс для сброса задач в хранилище
type Flusher interface {
	Flush(algorithms []algorithm.Algorithm) ([]algorithm.Algorithm, error)
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

func (f *flusher) Flush(algorithmsInput []algorithm.Algorithm) ([]algorithm.Algorithm, error) {
	bulks, err := algorithm.SplitAlgorithmsToBulks(algorithmsInput, f.chunkSize)
	if err != nil {
		return algorithmsInput, fmt.Errorf("cannot slit algorithms in bulks: %w", err)
	}

	if bulks == nil {
		return nil, nil
	}

	for len(bulks) > 0 {
		err := f.algorithmRepo.AddAlgorithms(bulks[0])
		if err != nil {
			return algorithm.AlgorithmBulksToSlice(bulks), fmt.Errorf("cannot same add algorithms into repo: %w", err)
		}
		bulks = bulks[1:]
	}

	return nil, nil
}
