package api_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"github.com/ozonva/ova-algorithm-api/internal/api"
	"github.com/ozonva/ova-algorithm-api/internal/mock_repo"
	desc "github.com/ozonva/ova-algorithm-api/pkg/ova-algorithm-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var _ = Describe("Api", func() {
	var (
		mockCtrl *gomock.Controller
		mockRepo *mock_repo.MockRepo
		s        desc.OvaAlgorithmApiServer
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mock_repo.NewMockRepo(mockCtrl)
		s = api.NewOvaAlgorithmApi(mockRepo)
	})

	AfterEach(func() {
		time.Sleep(100 * time.Millisecond)
		mockCtrl.Finish()
	})

	When("database create request is successful", func() {
		It("it should return nil error", func() {
			algo := algorithm.CreateSimpleAlgorithm(0)

			mockRepo.EXPECT().
				AddAlgorithms([]algorithm.Algorithm{algo}).
				Return(nil).
				Times(1)

			req := &desc.CreateAlgorithmRequestV1{
				Body: &desc.AlgorithmValueV1{
					Subject:     algo.Subject,
					Description: algo.Description,
				},
			}
			res, err := s.CreateAlgorithmV1(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
		})
	})

	When("database create request is failed", func() {
		It("it should return Unavailable error", func() {
			algo := algorithm.CreateSimpleAlgorithm(0)

			mockRepo.EXPECT().
				AddAlgorithms([]algorithm.Algorithm{algo}).
				Return(errors.New("cannot process sql")).
				Times(1)

			req := &desc.CreateAlgorithmRequestV1{
				Body: &desc.AlgorithmValueV1{
					Subject:     algo.Subject,
					Description: algo.Description,
				},
			}
			res, err := s.CreateAlgorithmV1(context.Background(), req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.Unavailable, "database store failed")))
			Expect(res).NotTo(BeNil())
		})
	})

	When("database describe can find entity", func() {
		It("it should return found algorithm and nil error", func() {
			algo := algorithm.CreateSimpleAlgorithm(0)

			const algorithmId = 0

			mockRepo.EXPECT().
				DescribeAlgorithm(uint64(algo.UserID)).
				Return(&algo, nil).
				Times(1)

			req := &desc.DescribeAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: algorithmId,
				},
			}

			res, err := s.DescribeAlgorithmV1(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Body.Id).To(Equal(int64(algo.UserID)))
			Expect(res.Body.Subject).To(Equal(algo.Subject))
			Expect(res.Body.Description).To(Equal(algo.Description))
		})
	})

	When("database describe cannot find entity", func() {
		It("it should return NotFound", func() {
			algo := algorithm.CreateSimpleAlgorithm(0)

			const algorithmId = 0

			mockRepo.EXPECT().
				DescribeAlgorithm(uint64(algo.UserID)).
				Return(nil, nil).
				Times(1)

			req := &desc.DescribeAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: algorithmId,
				},
			}

			res, err := s.DescribeAlgorithmV1(context.Background(), req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.NotFound, "identity not found")))
			Expect(res).To(BeNil())
		})
	})

	When("database describe got error while request", func() {
		It("it should return Unavailable", func() {
			algo := algorithm.CreateSimpleAlgorithm(0)

			const algorithmId = 0

			mockRepo.EXPECT().
				DescribeAlgorithm(uint64(algo.UserID)).
				Return(nil, errors.New("some error")).
				Times(1)

			req := &desc.DescribeAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: algorithmId,
				},
			}

			res, err := s.DescribeAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.Unavailable, "database fetch error")))
			Expect(res).To(BeNil())
		})
	})

	When("database list got error while request", func() {
		It("it should return Unavailable", func() {
			const limit = 3
			const offset = 5

			mockRepo.EXPECT().
				ListAlgorithms(uint64(limit), uint64(offset)).
				Return(nil, errors.New("some error")).
				Times(1)

			req := &desc.ListAlgorithmsRequestV1{
				Offset: &desc.AlgorithmIdV1{
					Id: offset,
				},
				Limit: 3,
			}

			res, err := s.ListAlgorithmsV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.Unavailable, "database list error")))
			Expect(res).To(BeNil())
		})
	})

	When("no entities have been found to be listed", func() {
		It("it should return NotFound", func() {
			const limit = 3
			const offset = 5

			mockRepo.EXPECT().
				ListAlgorithms(uint64(limit), uint64(offset)).
				Return(nil, nil).
				Times(1)

			req := &desc.ListAlgorithmsRequestV1{
				Offset: &desc.AlgorithmIdV1{
					Id: offset,
				},
				Limit: 3,
			}

			res, err := s.ListAlgorithmsV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.NotFound, "entities not found")))
			Expect(res).To(BeNil())
		})
	})

	When("two entities exists", func() {
		It("it should return two entities", func() {
			algos := algorithm.CreateSimpleAlgorithmListRangeInclusive(5, 6)
			const limit = 3
			const offset = 5

			mockRepo.EXPECT().
				ListAlgorithms(uint64(limit), uint64(offset)).
				Return(algos, nil).
				Times(1)

			req := &desc.ListAlgorithmsRequestV1{
				Offset: &desc.AlgorithmIdV1{
					Id: offset,
				},
				Limit: 3,
			}

			res, err := s.ListAlgorithmsV1(context.Background(), req)

			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())

			Expect(res.Body).NotTo(BeNil())
			Expect(len(res.Body)).To(Equal(len(algos)))

			for i := 0; i < len(algos); i++ {
				Expect(res.Body[i].Id).To(Equal(int64(algos[i].UserID)))
				Expect(res.Body[i].Subject).To(Equal(algos[i].Subject))
				Expect(res.Body[i].Description).To(Equal(algos[i].Description))
			}
		})
	})

	When("database remove got error while request", func() {
		It("it should return Unavailable", func() {
			const id = 3

			mockRepo.EXPECT().
				RemoveAlgorithm(uint64(id)).
				Return(false, errors.New("some error")).
				Times(1)

			req := &desc.RemoveAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: id,
				},
			}

			res, err := s.RemoveAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.Unavailable, "database delete error")))
			Expect(res).To(BeNil())
		})
	})

	When("no entities have been found to for removal", func() {
		It("it should return NotFound", func() {
			const id = 3

			mockRepo.EXPECT().
				RemoveAlgorithm(uint64(id)).
				Return(false, nil).
				Times(1)

			req := &desc.RemoveAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: id,
				},
			}

			res, err := s.RemoveAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.NotFound, "identity not found")))
			Expect(res).ToNot(BeNil())
		})
	})

	When("database remove can find entity", func() {
		It("it should return nil error", func() {
			const id = 3

			mockRepo.EXPECT().
				RemoveAlgorithm(uint64(id)).
				Return(true, nil).
				Times(1)

			req := &desc.RemoveAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: id,
				},
			}

			res, err := s.RemoveAlgorithmV1(context.Background(), req)

			Expect(err).NotTo(HaveOccurred())
			Expect(res).ToNot(BeNil())
		})
	})
})
