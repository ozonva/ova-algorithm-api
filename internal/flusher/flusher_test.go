package flusher_test

import (
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"
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
			It("should return nil", func() {
				notFlushed := flush.Flush(nil)
				Expect(notFlushed).To(BeNil())
			})
		})

		When("empty slice is used", func() {
			It("should return nil", func() {
				notFlushed := flush.Flush(nil)
				Expect(notFlushed).To(BeNil())
			})
		})

		When("no error on flush", func() {
			It("should return nil", func() {
				const listSize = 5
				mockRepo.EXPECT().
					AddAlgorithms(algorithm.CreateSimpleAlgorithmList(listSize)).
					Return(nil).
					Times(1)
				notFlushed := flush.Flush(algorithm.CreateSimpleAlgorithmList(listSize))
				Expect(notFlushed).To(BeNil())
			})
		})

		When("fails on flush", func() {
			It("should return the same slice", func() {
				const listSize = 5
				mockRepo.EXPECT().
					AddAlgorithms(algorithm.CreateSimpleAlgorithmList(listSize)).
					Return(errors.New("cannot save to repo")).
					Times(1)
				notFlushed := flush.Flush(algorithm.CreateSimpleAlgorithmList(listSize))
				Expect(notFlushed).To(Equal(algorithm.CreateSimpleAlgorithmList(listSize)))
			})
		})
	})

	Context("with bulkSize of 2)", func() {
		BeforeEach(func() {
			flush = flusher.NewFlusher(2, mockRepo)
		})

		It("accept nil slices", func() {
			notFlushed := flush.Flush(nil)
			Expect(notFlushed).To(BeNil())
		})

		It("accept empty slices", func() {
			notFlushed := flush.Flush(make([]algorithm.Algorithm, 0))
			Expect(notFlushed).To(BeNil())
		})

		Context("flushing with no errors", func() {
			When("list of one element", func() {
				It("accept all and return nil", func() {
					const listSize = 1
					mockRepo.EXPECT().
						AddAlgorithms(algorithm.CreateSimpleAlgorithmList(listSize)).
						Return(nil).
						Times(1)
					notFlushed := flush.Flush(algorithm.CreateSimpleAlgorithmList(listSize))
					Expect(notFlushed).To(BeNil())
				})
			})

			When("list two elements", func() {
				It("accept all and return nil", func() {
					const listSize = 2
					mockRepo.EXPECT().
						AddAlgorithms(algorithm.CreateSimpleAlgorithmList(listSize)).
						Return(nil).
						Times(1)
					notFlushed := flush.Flush(algorithm.CreateSimpleAlgorithmList(listSize))
					Expect(notFlushed).To(BeNil())
				})
			})

			When("list tree elements", func() {
				It("accept all and return nil", func() {
					gomock.InOrder(
						mockRepo.EXPECT().
							AddAlgorithms(algorithm.CreateSimpleAlgorithmListRangeInclusive(0, 1)).
							Return(nil).
							Times(1),
						mockRepo.EXPECT().
							AddAlgorithms(algorithm.CreateSimpleAlgorithmListRangeInclusive(2, 2)).
							Return(nil).
							Times(1),
					)
					notFlushed := flush.Flush(algorithm.CreateSimpleAlgorithmListRangeInclusive(0, 2))
					Expect(notFlushed).To(BeNil())
				})
			})
		})

		Context("flushing with all errors", func() {
			When("list of one element", func() {
				It("returns list of all input algorithms if cannot flush", func() {
					const listSize = 1
					list := algorithm.CreateSimpleAlgorithmList(listSize)
					mockRepo.EXPECT().
						AddAlgorithms(list).
						Return(fmt.Errorf("cannot flush to repo")).
						Times(1)
					notFlushed := flush.Flush(list)
					Expect(notFlushed).To(Equal(list))
				})
			})
		})

		Context("error on second flush of tree flushes", func() {
			When("list one of five element", func() {
				It("returns list of all input algorithms if cannot flush", func() {
					list := algorithm.CreateSimpleAlgorithmListRangeInclusive(0, 4)
					gomock.InOrder(
						mockRepo.EXPECT().
							AddAlgorithms(algorithm.CreateSimpleAlgorithmListRangeInclusive(0, 1)).
							Return(nil).
							Times(1),
						mockRepo.EXPECT().
							AddAlgorithms(algorithm.CreateSimpleAlgorithmListRangeInclusive(2, 3)).
							Return(fmt.Errorf("cannot flush to repo")).
							Times(1),
						mockRepo.EXPECT().
							AddAlgorithms(algorithm.CreateSimpleAlgorithmListRangeInclusive(4, 4)).
							Return(nil).
							Times(1),
					)
					notFlushed := flush.Flush(list)
					Expect(notFlushed).To(Equal(algorithm.CreateSimpleAlgorithmListRangeInclusive(2, 3)))
				})
			})
		})
	})
})
