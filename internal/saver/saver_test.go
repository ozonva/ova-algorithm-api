package saver_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
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
	})

	AfterEach(func() {
		time.Sleep(100 * time.Millisecond)
		mockCtrl.Finish()
	})

	Context("timer set to infinity, capacity 2", func() {
		BeforeEach(func() {
			s = saver.NewSaver(2, mocFlusher, time.Hour) // 1 hour = infinity
		})

		AfterEach(func() {
			s.Close()
		})

		When("no data added", func() {
			It("not call flusher at all", func() {
			})
		})

		When("only one entity has been added", func() {
			It("should only flush it once", func() {

				listOf1 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 1)

				mocFlusher.EXPECT().
					Flush(listOf1).
					Return(nil).
					Times(1)

				s.Save(algorithm.CreateSimpleAlgorithm(1))
			})
		})

		When("no data added", func() {
			It("not call flusher at all", func() {
			})
		})

		When("only two entity has been added", func() {
			It("should only flush it once", func() {

				listOf1_2 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 2)

				mocFlusher.EXPECT().
					Flush(listOf1_2).
					Return(nil).
					Times(1)

				s.Save(algorithm.CreateSimpleAlgorithm(1))
				s.Save(algorithm.CreateSimpleAlgorithm(2))
			})
		})

		When("only two entity has been added with error for both on first flush", func() {
			It("should only flush it once", func() {

				listOf1_2 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 2)

				mocFlusher.EXPECT().
					Flush(listOf1_2).
					Return(listOf1_2).
					Times(1)

				s.Save(algorithm.CreateSimpleAlgorithm(1))
				s.Save(algorithm.CreateSimpleAlgorithm(2))

				// flush on close
				mocFlusher.EXPECT().
					Flush(listOf1_2).
					Return(nil).
					Times(1)
			})
		})

		When("only two entity has been added with error for the first algo on first flush", func() {
			It("should only flush it once", func() {
				listOf1 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 1)
				listOf1_2 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 2)

				mocFlusher.EXPECT().
					Flush(listOf1_2).
					Return(listOf1).
					Times(1)

				s.Save(algorithm.CreateSimpleAlgorithm(1))
				s.Save(algorithm.CreateSimpleAlgorithm(2))

				// flush on close
				mocFlusher.EXPECT().
					Flush(listOf1).
					Return(nil).
					Times(1)
			})
		})

		When("only three entity has been added", func() {
			It("should only flush it twice", func() {

				listOf1_2 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 2)
				mocFlusher.EXPECT().
					Flush(listOf1_2).
					Return(nil).
					Times(1)

				s.Save(algorithm.CreateSimpleAlgorithm(1))
				s.Save(algorithm.CreateSimpleAlgorithm(2))

				time.Sleep(100 * time.Millisecond)

				s.Save(algorithm.CreateSimpleAlgorithm(3))

				listOf3 := algorithm.CreateSimpleAlgorithmListRangeInclusive(3, 3)
				mocFlusher.EXPECT().
					Flush(listOf3).
					Return(nil).
					Times(1)
			})
		})
	})

	Context("timer set to 2 second, capacity 2", func() {
		BeforeEach(func() {
			s = saver.NewSaver(2, mocFlusher, 2*time.Second)
		})

		AfterEach(func() {
			s.Close()
		})

		When("no data added", func() {
			It("not call flusher at all after 3 seconds", func() {
				time.Sleep(3 * time.Second)
			})
		})

		When("only one entity has been added", func() {
			It("should only flush it once after 3 seconds", func() {
				listOf1 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 1)

				mocFlusher.EXPECT().
					Flush(listOf1).
					Return(nil).
					Times(1)

				s.Save(algorithm.CreateSimpleAlgorithm(1))

				time.Sleep(3 * time.Second)
			})
		})

		When("two entities added, each flushed after timeout", func() {
			It("should only flush it once after 3 seconds", func() {
				listOf1 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 1)
				mocFlusher.EXPECT().
					Flush(listOf1).
					Return(nil).
					Times(1)

				s.Save(algorithm.CreateSimpleAlgorithm(1))
				time.Sleep(3 * time.Second)

				listOf2 := algorithm.CreateSimpleAlgorithmListRangeInclusive(2, 2)
				mocFlusher.EXPECT().
					Flush(listOf2).
					Return(nil).
					Times(1)

				s.Save(algorithm.CreateSimpleAlgorithm(2))
				time.Sleep(2 * time.Second)
			})
		})

		When("two entities added, each flushed after timeout, first flush fails", func() {
			It("should only flush it once after 3 seconds", func() {
				listOf1 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 1)
				mocFlusher.EXPECT().
					Flush(listOf1).
					Return(listOf1).
					Times(1)

				s.Save(algorithm.CreateSimpleAlgorithm(1))
				time.Sleep(3 * time.Second)

				listOf1_2 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 2)
				mocFlusher.EXPECT().
					Flush(listOf1_2).
					Return(nil).
					Times(1)

				s.Save(algorithm.CreateSimpleAlgorithm(2))
			})
		})
	})
})
