package algorithm

import (
	"fmt"
	"github.com/ozonva/ova-algorithm-api/internal/utils"
)

type Algorithm struct {
	UserID      uint64
	Subject     string
	Description string
}

// SplitAlgorithmsToBulks splits slice of []int into chunks of len chunkSize
func SplitAlgorithmsToBulks(algorithms []Algorithm, chunkSize int) ([][]Algorithm, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("chunkSize(%v) is negative or equal zero", chunkSize)
	}

	if len(algorithms) == 0 {
		return nil, nil
	}

	chunks := utils.CalculateChunks(len(algorithms), chunkSize)
	slices := make([][]Algorithm, chunks)

	for idx := 0; idx < chunks; idx++ {
		var chunkBegin = idx * chunkSize
		var chunkEnd = chunkBegin + chunkSize
		if chunkEnd > len(algorithms) {
			chunkEnd = len(algorithms)
		}
		slices[idx] = algorithms[chunkBegin:chunkEnd]
	}

	return slices, nil
}

// AlgorithmBulksToSlice concatenates bulks of Algorithm into
// single slice
func AlgorithmBulksToSlice(bulks [][]Algorithm) []Algorithm {
	flatSize := 0
	for i := 0; i < len(bulks); i++ {
		flatSize += len(bulks[i])
	}
	if flatSize == 0 {
		return nil
	}
	algorithms := make([]Algorithm, 0, flatSize)
	for i := 0; i < len(bulks); i++ {
		algorithms = append(algorithms, bulks[i]...)
	}
	return algorithms
}

// AlgorithmSliceToMap converts slice of Algorithm to map[uint64]Algorithm
// so UserId becomes key of a map. If duplicate UserId occurs in input slice
// error is returned
func AlgorithmSliceToMap(algorithms []Algorithm) (map[uint64]Algorithm, error) {
	if len(algorithms) == 0 {
		return make(map[uint64]Algorithm, 0), nil
	}

	resultMap := make(map[uint64]Algorithm, len(algorithms))

	for i := 0; i < len(algorithms); i++ {
		UserId := algorithms[i].UserID
		_, found := resultMap[UserId]
		if !found {
			resultMap[UserId] = algorithms[i]
		} else {
			return nil, fmt.Errorf("duplicate UserIDs: %v", UserId)
		}
	}

	return resultMap, nil
}
