package algorithm

import (
	"fmt"

	"github.com/ozonva/ova-algorithm-api/internal/numerics"
	"github.com/ozonva/ova-algorithm-api/internal/utils"
)

type Algorithm struct {
	UserID      uint64
	Subject     string
	Description string
}

// SplitAlgorithmsToBulks splits slice of []int into chunks of len chunkSize
func SplitAlgorithmsToBulks(algorithms []Algorithm, chunkSize uint) [][]Algorithm {
	const MaxInt = uint(numerics.MaxInt)
	if chunkSize > MaxInt {
		chunkSize = MaxInt
	}
	chunkSizeInt := int(chunkSize)

	if len(algorithms) == 0 {
		return nil
	}

	if chunkSizeInt == 0 {
		return [][]Algorithm{algorithms}
	}

	chunks := utils.CalculateChunks(len(algorithms), chunkSizeInt)
	slices := make([][]Algorithm, chunks)

	for idx := 0; idx < chunks; idx++ {
		var chunkBegin = idx * chunkSizeInt
		var chunkEnd = chunkBegin + chunkSizeInt
		if chunkEnd > len(algorithms) {
			chunkEnd = len(algorithms)
		}
		slices[idx] = algorithms[chunkBegin:chunkEnd]
	}

	return slices
}

// AlgorithmSliceToMap converts slice of Algorithm to map[uint64]Algorithm
// so UserId becomes key of a map. If duplicate UserId occurs in input slice
// error is returned
func AlgorithmSliceToMap(algorithms []Algorithm) (map[uint64]Algorithm, error) {
	if len(algorithms) == 0 {
		return make(map[uint64]Algorithm), nil
	}

	resultMap := make(map[uint64]Algorithm, len(algorithms))

	for i := 0; i < len(algorithms); i++ {
		id := algorithms[i].UserID
		_, found := resultMap[id]
		if !found {
			resultMap[id] = algorithms[i]
		} else {
			return nil, fmt.Errorf("duplicate UserIDs: %v", id)
		}
	}

	return resultMap, nil
}
