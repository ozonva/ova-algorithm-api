package saver_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/ozonva/ova-algorithm-api/internal/mock_flusher"
	saver "github.com/ozonva/ova-algorithm-api/internal/saver"
	"time"
)

var _ = Describe("Saver", func() {
	var (
		mockCtrl   *gomock.Controller
		mocFlusher *mock_flusher.MockFlusher
		s          saver.Saver
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mocFlusher = mock_flusher.NewMockFlusher(mockCtrl)
		s = saver.NewSaver(2, mocFlusher, time.Hour)
	})

	AfterEach(func() {
		s.Close()
	})
})
