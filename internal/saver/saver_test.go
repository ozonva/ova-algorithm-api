package saver_test

import (
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"github.com/ozonva/ova-algorithm-api/internal/mock_flusher"
	"github.com/ozonva/ova-algorithm-api/internal/saver"
)

var _ = Describe("Saver", func() {
	var (
		mockCtrl   *gomock.Controller
		mocFlusher *mock_flusher.MockFlusher
		s          saver.Saver
		startTime  time.Time
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mocFlusher = mock_flusher.NewMockFlusher(mockCtrl)
		startTime = time.Now()
	})

	AfterEach(func() {
		time.Sleep(100 * time.Millisecond)
		mockCtrl.Finish()
	})

	Context("timer set to 2 second, capacity 2", func() {
		const Capacity = 2
		const FlushPeriod = 2 * time.Second

		BeforeEach(func() {
			s = saver.NewSaver(Capacity, mocFlusher, FlushPeriod)
		})

		AfterEach(func() {
			s.Stop()
		})

		When("no data is added to the store", func() {
			BeforeEach(func() {
				By("initially asserting the zero-sized storage with configured capacity")
				length, capacity := s.GetLenAndCap()
				Expect(length).To(Equal(0))
				Expect(capacity).To(Equal(Capacity))
			})

			AfterEach(func() {
				By("finally asserting the zero-sized storage with configured capacity")
				length, capacity := s.GetLenAndCap()
				Expect(length).To(Equal(0))
				Expect(capacity).To(Equal(Capacity))
			})

			It("shall not call flush after 1.5 flush period", func() {
				sleepUntil(startTime.Add(FlushPeriod * 3 / 2))
			})

			It("shall not call flush after explicit Close call", func() {
				s.Close()
			})
		})

		When("one entity has been added to the storage", func() {
			BeforeEach(func() {
				By("initially asserting the zero-sized storage with configured capacity", func() {
					length, capacity := s.GetLenAndCap()
					Expect(length).To(Equal(0))
					Expect(capacity).To(Equal(Capacity))
				})

				By("saving algorithm[1] to the storage", func() {
					err := s.Save(algorithm.CreateSimpleAlgorithm(1))
					Expect(err).ShouldNot(HaveOccurred())
				})

				By("asserting the one-sized storage with configured capacity", func() {
					length, capacity := s.GetLenAndCap()
					Expect(length).To(Equal(1))
					Expect(capacity).To(Equal(Capacity))
				})
			})

			AfterEach(func() {
				By("finally asserting the zero-sized storage with configured capacity")
				length, capacity := s.GetLenAndCap()
				Expect(length).To(Equal(0))
				Expect(capacity).To(Equal(Capacity))
			})

			It("shall flush entities on timeout", func() {
				By("configuring flush mock to flush all on list of 1", func() {
					listOf1 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 1)
					mocFlusher.EXPECT().
						Flush(listOf1).
						Return(nil).
						Times(1)
				})

				By("waiting storage has flushed after 1.5 flush period", func() {
					sleepUntil(startTime.Add(FlushPeriod * 3 / 2))
				})
			})

			It("shall flush entities on explicit Close call", func() {
				By("configuring flush mock to flush all on list of 1", func() {
					listOf1 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 1)
					mocFlusher.EXPECT().
						Flush(listOf1).
						Return(nil).
						Times(1)
				})

				By("explicitly calling s.Close()", func() {
					s.Close()
				})
			})

			It("shall flush entities if the storage is full", func() {
				By("configuring flush mock to flush all on list of 1,2", func() {
					listOf12 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 2)
					mocFlusher.EXPECT().
						Flush(listOf12).
						Return(nil).
						Times(1)
				})

				By("saving 2 and checking it's result", func() {
					err := s.Save(algorithm.CreateSimpleAlgorithm(2))
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		When("storage is full", func() {
			BeforeEach(func() {
				By("initially asserting the zero-sized storage with configured capacity", func() {
					length, capacity := s.GetLenAndCap()
					Expect(length).To(Equal(0))
					Expect(capacity).To(Equal(Capacity))
				})

				By("saving algorithm[1] to the storage", func() {
					err := s.Save(algorithm.CreateSimpleAlgorithm(1))
					Expect(err).ShouldNot(HaveOccurred())
				})

				By("asserting the one-sized storage with configured capacity", func() {
					length, capacity := s.GetLenAndCap()
					Expect(length).To(Equal(1))
					Expect(capacity).To(Equal(Capacity))
				})

				By("configuring flush mock to fail flush on list of 1,2", func() {
					listOf12 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 2)
					mocFlusher.EXPECT().
						Flush(listOf12).
						Return(listOf12).
						Times(1)
				})

				By("saving algorithm[2] to the storage", func() {
					err := s.Save(algorithm.CreateSimpleAlgorithm(2))
					Expect(err).ShouldNot(HaveOccurred())
				})

				By("asserting the one-sized storage with configured capacity", func() {
					length, capacity := s.GetLenAndCap()
					Expect(length).To(Equal(2))
					Expect(capacity).To(Equal(Capacity))
				})
			})

			It("shall fail Save to the storage", func() {
				By("calling Save with algorithm[3]", func() {
					err := s.Save(algorithm.CreateSimpleAlgorithm(3))
					Expect(err).Should(HaveOccurred())
				})

				By("asserting storage is still full", func() {
					length, capacity := s.GetLenAndCap()
					Expect(length).To(Equal(2))
					Expect(capacity).To(Equal(Capacity))
				})
			})

			It("shall flush the storage after 1.5 flush period", func() {
				By("configuring flush mock to flush all on list of 1,2", func() {
					listOf12 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 2)
					mocFlusher.EXPECT().
						Flush(listOf12).
						Return(nil).
						Times(1)
				})

				By("waiting storage has flushed after 1.5 flush period", func() {
					sleepUntil(startTime.Add(FlushPeriod * 3 / 2))
				})

				By("asserting storage is empty", func() {
					length, capacity := s.GetLenAndCap()
					Expect(length).To(Equal(0))
					Expect(capacity).To(Equal(Capacity))
				})
			})
		})
	})
})

// added to mitigate time deviations
func sleepUntil(deadline time.Time) {
	elapsed := time.Until(deadline)
	if elapsed < 0 {
		panic("time expectations broken")
	}
	time.Sleep(elapsed)
}
