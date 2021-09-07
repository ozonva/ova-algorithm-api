package saver

import (
	"errors"
	"sync"
	"time"

	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"github.com/ozonva/ova-algorithm-api/internal/flusher"
)

type Saver interface {
	// Save add Algorithm to the Saver
	Save(algorithm algorithm.Algorithm) error
	// GetLenAndCap returns len and cap of internal storage
	GetLenAndCap() (int, int)
	// Close forces immediate flush
	Close()
	// Stop stops periodic saves
	Stop()
}

type saver struct {
	ticker *time.Ticker
	sync.Mutex
	store   []algorithm.Algorithm
	flusher flusher.Flusher
}

// NewSaver возвращает Saver с поддержкой переодического сохранения
func NewSaver(
	capacity uint,
	flusher flusher.Flusher,
	duration time.Duration,

) Saver {
	saver := &saver{
		store:   make([]algorithm.Algorithm, 0, capacity),
		flusher: flusher,
		ticker:  time.NewTicker(duration),
	}
	go func() {
		for range saver.ticker.C {
			saver.Close()
		}
	}()
	return saver
}

func (s *saver) Stop() {
	s.ticker.Stop()
}

func (s *saver) Save(algorithm algorithm.Algorithm) error {
	s.Lock()
	defer s.Unlock()
	if len(s.store) == cap(s.store) {
		return errors.New("cannot save, storage is full")
	}
	s.pushAlgorithm(algorithm)
	return nil
}

func (s *saver) Close() {
	s.Lock()
	defer s.Unlock()
	if len(s.store) > 0 {
		s.flushStore()
	}
}

func (s *saver) GetLenAndCap() (int, int) {
	s.Lock()
	defer s.Unlock()
	return len(s.store), cap(s.store)
}

func (s *saver) pushAlgorithm(entity algorithm.Algorithm) {
	s.store = append(s.store, entity)

	if len(s.store) == cap(s.store) {
		s.flushStore()
	}
}

func (s *saver) flushStore() {
	capacity := cap(s.store)
	left := s.flusher.Flush(s.store)

	if cap(left) != capacity {
		s.store = make([]algorithm.Algorithm, 0, capacity)
		s.store = append(s.store, left...)
	} else {
		s.store = left
	}
}
