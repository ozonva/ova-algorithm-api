package repo

import (
	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
)

type Repo interface {
	AddAlgorithms(algorithm []algorithm.Algorithm) error
	ListAlgorithms(limit, offset uint64) ([]algorithm.Algorithm, error)
	DescribeAlgorithm(algorithmID uint64) (*algorithm.Algorithm, error)
}
