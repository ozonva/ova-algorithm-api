package repo

import (
	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
)

type Repo interface {
	AddAlgorithms(algorithm []algorithm.Algorithm) error
	ListAlgorithms(limit, offset uint64) ([]algorithm.Algorithm, error)
	DescribeAlgorithm(algorithmID uint64) (*algorithm.Algorithm, error)
}

func NewRepo() Repo {
	return &repo{}
}

type repo struct {

}

func (r *repo) AddAlgorithms(algorithm []algorithm.Algorithm) error {
	return nil
}

func (r *repo) ListAlgorithms(limit, offset uint64) ([]algorithm.Algorithm, error) {
	return nil, nil
}

func (r *repo) DescribeAlgorithm(algorithmID uint64) (*algorithm.Algorithm, error) {
	return nil, nil
}