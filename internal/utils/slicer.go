package utils

type SplitToChunksError int

const (
	NoError SplitToChunksError = iota
	ZeroOrNegativeChunksSize
)

func CalculateChunks(sliceSize int, chunkSize int) int {
	quotient, remainder := sliceSize/chunkSize, sliceSize%chunkSize

	chunks := quotient
	if remainder > 0 {
		chunks += 1
	}
	return chunks
}

// SplitToChunksInt splits slice of []int into chunks of len chunkSize.
// For invalid number of chunks the value ZeroOrNegativeChunksSize of
// SplitToChunksError is returned
func SplitToChunksInt(slice []int, chunkSize int) ([][]int, SplitToChunksError) {
	if chunkSize <= 0 {
		return nil, ZeroOrNegativeChunksSize
	}

	if len(slice) == 0 {
		return nil, NoError
	}

	chunks := CalculateChunks(len(slice), chunkSize)
	slices := make([][]int, chunks)

	for idx := 0; idx < chunks; idx++ {
		var chunkBegin = idx * chunkSize
		var chunkEnd = chunkBegin + chunkSize
		if chunkEnd > len(slice) {
			chunkEnd = len(slice)
		}
		slices[idx] = slice[chunkBegin:chunkEnd]
	}

	return slices, NoError
}
