package saver_test

import (
	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

		AfterEach(func() {
			s.Stop()
		})

		AssertSizeWithConfiguredCapacity := func(size int) {
			length, capacity := s.GetLenAndCap()
			Expect(length).To(Equal(size))
			Expect(capacity).To(Equal(Capacity))
		}

		When("no data is added to the store", func() {
			It("shall not call flush after 1.5 flush period", func() {
				By("starting saver", func() {
					s = saver.NewSaver(Capacity, mocFlusher, FlushPeriod)
				})

				By("initially asserting the zero-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(0)
				})

				By("waiting storage has flushed after 1.5 flush period", func() {
					sleepUntil(startTime.Add(FlushPeriod * 3 / 2))
				})

				By("finally asserting the zero-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(0)
				})
			})

			It("shall not call flush after explicit Close call", func() {
				By("starting saver", func() {
					s = saver.NewSaver(Capacity, mocFlusher, FlushPeriod)
				})

				By("initially asserting the zero-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(0)
				})

				By("explicitly calling close", func() {
					s.Close()
				})

				By("finally asserting the zero-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(0)
				})
			})
		})

		When("one entity has been added to the storage", func() {
			AssertOneAdded := func() {
				By("starting saver", func() {
					s = saver.NewSaver(Capacity, mocFlusher, FlushPeriod)
				})

				By("initially asserting the zero-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(0)
				})

				By("saving algorithm[1] to the storage", func() {
					err := s.Save(algorithm.CreateSimpleAlgorithm(1))
					Expect(err).ShouldNot(HaveOccurred())
				})

				By("asserting the one-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(1)
				})
			}

			It("shall flush entities on timeout", func() {
				By("configuring flush mock to flush all on list of 1", func() {
					listOf1 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 1)
					mocFlusher.EXPECT().
						Flush(listOf1).
						Return(nil).
						Times(1)
				})

				AssertOneAdded()

				By("waiting storage has flushed after 1.5 flush period", func() {
					sleepUntil(startTime.Add(FlushPeriod * 3 / 2))
				})

				By("finally asserting the zero-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(0)
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

				AssertOneAdded()

				By("explicitly calling s.Close()", func() {
					s.Close()
				})

				By("finally asserting the zero-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(0)
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

				AssertOneAdded()

				By("saving 2 and checking it's result", func() {
					err := s.Save(algorithm.CreateSimpleAlgorithm(2))
					Expect(err).ShouldNot(HaveOccurred())
				})

				By("asserting the zero-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(0)
				})
			})
		})

		When("storage is full", func() {

			AssertStorageIsFull := func() {
				By("starting saver", func() {
					s = saver.NewSaver(Capacity, mocFlusher, FlushPeriod)
				})

				By("initially asserting the zero-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(0)
				})

				By("saving algorithm[1] to the storage", func() {
					err := s.Save(algorithm.CreateSimpleAlgorithm(1))
					Expect(err).ShouldNot(HaveOccurred())
				})

				By("asserting the one-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(1)
				})

				By("saving algorithm[2] to the storage", func() {
					err := s.Save(algorithm.CreateSimpleAlgorithm(2))
					Expect(err).ShouldNot(HaveOccurred())
				})

				By("asserting the two-sized storage with configured capacity", func() {
					AssertSizeWithConfiguredCapacity(2)
				})
			}

			It("shall fail Save to the storage", func() {
				By("configuring flush mock to fail flush on list of 1,2", func() {
					listOf12 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 2)
					mocFlusher.EXPECT().
						Flush(listOf12).
						Return(listOf12).
						Times(1)
				})

				AssertStorageIsFull()

				By("calling Save with algorithm[3]", func() {
					err := s.Save(algorithm.CreateSimpleAlgorithm(3))
					Expect(err).Should(HaveOccurred())
				})

				By("asserting storage is still full", func() {
					AssertSizeWithConfiguredCapacity(Capacity)
				})
			})

			It("shall flush the storage after 1.5 flush period", func() {
				By("configuring flush mock to fail flush all on list of 1,2 on the first attempt and"+
					"succeed on the second", func() {
					listOf12 := algorithm.CreateSimpleAlgorithmListRangeInclusive(1, 2)
					gomock.InOrder(
						mocFlusher.EXPECT().
							Flush(listOf12).
							Return(listOf12).
							Times(1),
						mocFlusher.EXPECT().
							Flush(listOf12).
							Return(nil).
							Times(1),
					)
				})

				AssertStorageIsFull()

				By("failing algorithm[3] to the storage", func() {
					err := s.Save(algorithm.CreateSimpleAlgorithm(3))
					Expect(err).Should(HaveOccurred())
				})

				By("waiting storage has flushed after 1.5 flush period", func() {
					sleepUntil(startTime.Add(FlushPeriod * 3 / 2))
				})

				By("asserting storage is empty", func() {
					AssertSizeWithConfiguredCapacity(0)
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
