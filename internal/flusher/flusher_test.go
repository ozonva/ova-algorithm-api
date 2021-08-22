package flusher_test

import (
	"fmt"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"github.com/ozonva/ova-algorithm-api/internal/flusher"
	"github.com/ozonva/ova-algorithm-api/internal/mock_repo"
)

var _ = Describe("Flusher", func() {
	var (
		mockCtrl *gomock.Controller
		mockRepo *mock_repo.MockRepo
		flush    flusher.Flusher
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mock_repo.NewMockRepo(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("with zero bulk size", func() {
		BeforeEach(func() {
			flush = flusher.NewFlusher(0, mockRepo)
		})

		When("nil slice is used", func() {
			It("should return nil and error", func() {
				notFlushed, err := flush.Flush(nil)
				Expect(notFlushed).To(BeNil())
				Expect(err).NotTo(BeNil())
			})
		})

		When("empty slice is used", func() {
			It("should return nil and error", func() {
				notFlushed, err := flush.Flush(nil)
				Expect(notFlushed).To(BeNil())
				Expect(err).NotTo(BeNil())
			})
		})

		When("non empty slice is used", func() {
			It("should return input slice and error", func() {
				const listSize = 5
				notFlushed, err := flush.Flush(createSimpleAlgorithmList(listSize))
				Expect(notFlushed).To(Equal(createSimpleAlgorithmList(listSize)))
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Context("with bulkSize of 2)", func() {
		BeforeEach(func() {
			flush = flusher.NewFlusher(2, mockRepo)
		})

		It("accept nil slices", func() {
			notFlushed, err := flush.Flush(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(notFlushed).To(BeNil())
		})

		It("accept empty slices", func() {
			notFlushed, err := flush.Flush(make([]algorithm.Algorithm, 0))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(notFlushed).To(BeNil())
		})

		Context("flushing with no errors", func() {
			When("list of one element", func() {
				It("accept all and return nil", func() {
					const listSize = 1
					mockRepo.EXPECT().AddAlgorithms(createSimpleAlgorithmList(listSize)).Return(nil).Times(1)
					notFlushed, err := flush.Flush(createSimpleAlgorithmList(listSize))
					Expect(err).ShouldNot(HaveOccurred())
					Expect(notFlushed).To(BeNil())
				})
			})

			When("list two elements", func() {
				It("accept all and return nil", func() {
					const listSize = 2
					mockRepo.EXPECT().AddAlgorithms(createSimpleAlgorithmList(listSize)).Return(nil).Times(1)
					notFlushed, err := flush.Flush(createSimpleAlgorithmList(listSize))
					Expect(err).ShouldNot(HaveOccurred())
					Expect(notFlushed).To(BeNil())
				})
			})

			When("list tree elements", func() {
				It("accept all and return nil", func() {
					gomock.InOrder(
						mockRepo.EXPECT().AddAlgorithms(createSimpleAlgorithmListRangeInclusive(0, 1)).Return(nil).Times(1),
						mockRepo.EXPECT().AddAlgorithms(createSimpleAlgorithmListRangeInclusive(2, 2)).Return(nil).Times(1),
					)
					notFlushed, err := flush.Flush(createSimpleAlgorithmListRangeInclusive(0, 2))
					Expect(err).ShouldNot(HaveOccurred())
					Expect(notFlushed).To(BeNil())
				})
			})
		})

		Context("flushing with all errors", func() {
			When("list of one element", func() {
				It("returns list of all input algorithms if cannot flush", func() {
					const listSize = 1
					list := createSimpleAlgorithmList(listSize)
					mockRepo.EXPECT().AddAlgorithms(list).Return(fmt.Errorf("cannot flush to repo")).Times(1)
					notFlushed, err := flush.Flush(list)
					Expect(err).Should(HaveOccurred())
					Expect(notFlushed).To(Equal(list))
				})
			})
		})

		Context("error on second flush", func() {
			When("list one of three element", func() {
				It("returns list of all input algorithms if cannot flush", func() {
					const listSize = 1
					list := createSimpleAlgorithmListRangeInclusive(0, 2)
					gomock.InOrder(
						mockRepo.EXPECT().AddAlgorithms(createSimpleAlgorithmListRangeInclusive(0, 1)).Return(nil).Times(1),
						mockRepo.EXPECT().AddAlgorithms(createSimpleAlgorithmListRangeInclusive(2, 2)).Return(fmt.Errorf("cannot flush to repo")).Times(1),
					)
					notFlushed, err := flush.Flush(list)
					Expect(err).Should(HaveOccurred())
					Expect(notFlushed).To(Equal(createSimpleAlgorithmListRangeInclusive(2, 2)))
				})
			})
		})
	})
})

func createSimpleAlgorithmListRangeInclusive(begin, end int) []algorithm.Algorithm {
	if end < begin {
		panic(fmt.Sprintf("end(%v) should not less begin(%v)", end, begin))
	}
	size := end - begin + 1
	list := make([]algorithm.Algorithm, 0, size)
	for i := begin; i <= end; i++ {
		list = append(list, algorithm.Algorithm{
			UserID:      uint64(i),
			Subject:     fmt.Sprintf("Subject%v", i),
			Description: fmt.Sprintf("Description%v", i),
		})
	}
	return list
}

func createSimpleAlgorithmList(values int) []algorithm.Algorithm {
	return createSimpleAlgorithmListRangeInclusive(1, values)
}
