package saver

import (
	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"github.com/ozonva/ova-algorithm-api/internal/flusher"
	"time"
)

type Saver interface {
	Save(algorithm algorithm.Algorithm) // заменить на свою сущность
	Close()
}

type saver struct {
	ch chan<- algorithm.Algorithm
}

type saverRunner struct {
	maxCapacity uint
	ch          <-chan algorithm.Algorithm
	store       []algorithm.Algorithm
	flusher     flusher.Flusher
}

// NewSaver возвращает Saver с поддержкой переодического сохранения
func NewSaver(
	capacity uint,
	flusher flusher.Flusher,
	duration time.Duration,
) Saver {
	ch := make(chan algorithm.Algorithm)
	saver := &saver{
		ch: ch,
	}
	saverRunner := &saverRunner{
		maxCapacity: capacity,
		ch:          ch,
		store:       make([]algorithm.Algorithm, 0, capacity),
		flusher:     flusher,
	}
	go saverRunner.run(duration)
	return saver
}

func (s *saver) Save(algorithm algorithm.Algorithm) {
	s.ch <- algorithm
}

func (s *saver) Close() {
	close(s.ch)
}

func (r *saverRunner) run(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case entity, ok := <-r.ch:
			if !ok {
				r.flushStore()
				return
			}
			r.pushAlgorithm(entity)
		case <-ticker.C:
			r.flushStore()
		}
	}
}

func (r *saverRunner) pushAlgorithm(algorithm algorithm.Algorithm) {
	r.store = append(r.store, algorithm)
	if uint(len(r.store)) >= r.maxCapacity {
		r.flushStore()
	}
}

func (r *saverRunner) flushStore() {
	if len(r.store) > 0 {
		r.store = r.flusher.Flush(r.store)
	}
}
