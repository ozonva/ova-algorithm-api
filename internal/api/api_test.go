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
			algo := algorithm.CreateSimpleAlgorithm(1)

			const algorithmId = 1

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
			algo := algorithm.CreateSimpleAlgorithm(1)

			const algorithmId = 1

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
			algo := algorithm.CreateSimpleAlgorithm(1)

			const algorithmId = 1

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

	When("database describe go id(0) out of range", func() {
		It("it should return OutOfRange", func() {
			const algorithmId = 0

			req := &desc.DescribeAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: algorithmId,
				},
			}

			res, err := s.DescribeAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.OutOfRange, "id (0) is out of range 1 - 2,147,483,647")))
			Expect(res).To(BeNil())
		})
	})

	When("database describe go id(2147483648) out of range", func() {
		It("it should return OutOfRange", func() {
			const algorithmId = 2147483648

			req := &desc.DescribeAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: algorithmId,
				},
			}

			res, err := s.DescribeAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.OutOfRange, "id (2147483648) is out of range 1 - 2,147,483,647")))
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

	When("offset is out of range (-1) is provided", func() {
		It("it should return NotFound", func() {
			const limit = 3
			const offset = -1

			req := &desc.ListAlgorithmsRequestV1{
				Offset: &desc.AlgorithmIdV1{
					Id: offset,
				},
				Limit: limit,
			}

			res, err := s.ListAlgorithmsV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.OutOfRange, "offset (-1) is out of range 0 - 2,147,483,647")))
			Expect(res).To(BeNil())
		})
	})

	When("offset is out of range (2,147,483,648) is provided", func() {
		It("it should return NotFound", func() {
			const limit = 3
			const offset = 2147483648

			req := &desc.ListAlgorithmsRequestV1{
				Offset: &desc.AlgorithmIdV1{
					Id: offset,
				},
				Limit: limit,
			}

			res, err := s.ListAlgorithmsV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.OutOfRange, "offset (2147483648) is out of range 0 - 2,147,483,647")))
			Expect(res).To(BeNil())
		})
	})

	When("limit is out of range (0) is provided", func() {
		It("it should return NotFound", func() {
			const limit = 0
			const offset = 5

			req := &desc.ListAlgorithmsRequestV1{
				Offset: &desc.AlgorithmIdV1{
					Id: offset,
				},
				Limit: limit,
			}

			res, err := s.ListAlgorithmsV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.OutOfRange, "limit (0) is out of range 1 - 2,147,483,647")))
			Expect(res).To(BeNil())
		})
	})

	When("limit is out of range (2,147,483,648) is provided", func() {
		It("it should return NotFound", func() {
			const limit = 2147483648
			const offset = 5

			req := &desc.ListAlgorithmsRequestV1{
				Offset: &desc.AlgorithmIdV1{
					Id: offset,
				},
				Limit: limit,
			}

			res, err := s.ListAlgorithmsV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.OutOfRange, "limit (2147483648) is out of range 1 - 2,147,483,647")))
			Expect(res).To(BeNil())
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
			Expect(res).NotTo(BeNil())
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

	When("id out of range (0) provided for removal", func() {
		It("should return OutOfRange", func() {
			const id = 0

			req := &desc.RemoveAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: id,
				},
			}

			res, err := s.RemoveAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.OutOfRange, "id (0) is out of range 1 - 2,147,483,647")))
			Expect(res).ToNot(BeNil())
		})
	})

	When("id out of range (2,147,483,648) provided for removal", func() {
		It("should return OutOfRange", func() {
			const id = 2147483648

			req := &desc.RemoveAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: id,
				},
			}

			res, err := s.RemoveAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.OutOfRange, "id (2147483648) is out of range 1 - 2,147,483,647")))
			Expect(res).ToNot(BeNil())
		})
	})

	When("entity update got database error while request", func() {
		It("it should return Unavailable", func() {
			algo := algorithm.CreateSimpleAlgorithm(1)

			mockRepo.EXPECT().
				UpdateAlgorithm(algo).
				Return(false, errors.New("some error")).
				Times(1)

			req := &desc.UpdateAlgorithmRequestV1{
				Body: &desc.AlgorithmV1{
					Id:          int64(algo.UserID),
					Subject:     algo.Subject,
					Description: algo.Description,
				},
			}

			res, err := s.UpdateAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.Unavailable, "database update error")))
			Expect(res).NotTo(BeNil())
		})
	})

	When("no entity has been found to for update", func() {
		It("it should return NotFound", func() {
			algo := algorithm.CreateSimpleAlgorithm(1)

			mockRepo.EXPECT().
				UpdateAlgorithm(algo).
				Return(false, nil).
				Times(1)

			req := &desc.UpdateAlgorithmRequestV1{
				Body: &desc.AlgorithmV1{
					Id:          int64(algo.UserID),
					Subject:     algo.Subject,
					Description: algo.Description,
				},
			}

			res, err := s.UpdateAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.NotFound, "identity not found")))
			Expect(res).ToNot(BeNil())
		})
	})

	When("update were able to find entity", func() {
		It("it should return nil error", func() {
			algo := algorithm.CreateSimpleAlgorithm(1)

			mockRepo.EXPECT().
				UpdateAlgorithm(algo).
				Return(true, nil).
				Times(1)

			req := &desc.UpdateAlgorithmRequestV1{
				Body: &desc.AlgorithmV1{
					Id:          int64(algo.UserID),
					Subject:     algo.Subject,
					Description: algo.Description,
				},
			}

			res, err := s.UpdateAlgorithmV1(context.Background(), req)

			Expect(err).NotTo(HaveOccurred())
			Expect(res).ToNot(BeNil())
		})
	})

	When("id out of range (0) provided for update", func() {
		It("should return OutOfRange", func() {
			algo := algorithm.CreateSimpleAlgorithm(0)

			req := &desc.UpdateAlgorithmRequestV1{
				Body: &desc.AlgorithmV1{
					Id:          int64(algo.UserID),
					Subject:     algo.Subject,
					Description: algo.Description,
				},
			}

			res, err := s.UpdateAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.OutOfRange, "id (0) is out of range 1 - 2,147,483,647")))
			Expect(res).ToNot(BeNil())
		})
	})

	When("id out of range (2,147,483,648) provided for update", func() {
		It("should return OutOfRange", func() {
			algo := algorithm.CreateSimpleAlgorithm(2147483648)

			req := &desc.UpdateAlgorithmRequestV1{
				Body: &desc.AlgorithmV1{
					Id:          int64(algo.UserID),
					Subject:     algo.Subject,
					Description: algo.Description,
				},
			}

			res, err := s.UpdateAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.OutOfRange, "id (2147483648) is out of range 1 - 2,147,483,647")))
			Expect(res).ToNot(BeNil())
		})
	})
})
