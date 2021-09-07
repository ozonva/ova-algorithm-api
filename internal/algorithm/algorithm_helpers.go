package algorithm

import (
	"fmt"
)

// CreateSimpleAlgorithm generates an algorithm for test purposes.
// All fields are generated as <filed_name><id>. The id is assigned
// to itself
func CreateSimpleAlgorithm(idx int) Algorithm {
	return Algorithm{
		UserID:      uint64(idx),
		Subject:     fmt.Sprintf("Subject%v", idx),
		Description: fmt.Sprintf("Description%v", idx),
	}
}

// CreateSimpleAlgorithmListRangeInclusive creates as closed range of
// algorithms using CreateSimpleAlgorithm function
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

// CreateSimpleAlgorithmList creates a list of algorithms with provided size
// Values indexes stars with 1 and reaches provided number of values
func CreateSimpleAlgorithmList(values int) []Algorithm {
	return CreateSimpleAlgorithmListRangeInclusive(1, values)
}
