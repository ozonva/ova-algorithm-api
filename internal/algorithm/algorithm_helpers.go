package algorithm

import (
	"fmt"
)

func CreateSimpleAlgorithm(idx int) Algorithm {
	return Algorithm{
		UserID:      uint64(idx),
		Subject:     fmt.Sprintf("Subject%v", idx),
		Description: fmt.Sprintf("Description%v", idx),
	}
}

func CreateSimpleAlgorithmListRangeInclusive(begin, end int) []Algorithm {
	if end < begin {
		panic(fmt.Sprintf("end(%v) should not less begin(%v)", end, begin))
	}
	size := end - begin + 1
	list := make([]Algorithm, 0, size)
	for i := begin; i <= end; i++ {
		list = append(list, CreateSimpleAlgorithm(i))
	}
	return list
}

func CreateSimpleAlgorithmList(values int) []Algorithm {
	return CreateSimpleAlgorithmListRangeInclusive(1, values)
}
