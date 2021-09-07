package api_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	saramaMocks "github.com/Shopify/sarama/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"github.com/ozonva/ova-algorithm-api/internal/api"
	"github.com/ozonva/ova-algorithm-api/internal/mock_repo"
	"github.com/ozonva/ova-algorithm-api/internal/notification"
	desc "github.com/ozonva/ova-algorithm-api/pkg/ova-algorithm-api"
)

var _ = Describe("Api", func() {
	var (
		mockCtrl   *gomock.Controller
		mockRepo   *mock_repo.MockRepo
		notifyMock *saramaMocks.AsyncProducer
		s          desc.OvaAlgorithmApiServer
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mock_repo.NewMockRepo(mockCtrl)
		notifyMock = saramaMocks.NewAsyncProducer(GinkgoT(), nil)
		//notifyMock.ExpectInputAndSucceed()
		s = api.NewOvaAlgorithmAPI(mockRepo, notifyMock)
	})

	AfterEach(func() {
		notifyMock.Close()
		mockCtrl.Finish()
	})

	When("database create request is successful", func() {
		It("it should return nil error", func() {
			algo := algorithm.CreateSimpleAlgorithm(0)

			notifyMock.ExpectInputWithCheckerFunctionAndSucceed(
				createAlgorithmNotificationChecker(0, notification.OpCreate))

			mockRepo.EXPECT().
				AddAlgorithms(context.Background(), []algorithm.Algorithm{algo}).
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
				AddAlgorithms(context.Background(), []algorithm.Algorithm{algo}).
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
			const id = 1

			algo := algorithm.CreateSimpleAlgorithm(id)

			mockRepo.EXPECT().
				DescribeAlgorithm(context.Background(), algo.UserID).
				Return(&algo, nil).
				Times(1)

			req := &desc.DescribeAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: id,
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
			const id = 1

			algo := algorithm.CreateSimpleAlgorithm(id)

			mockRepo.EXPECT().
				DescribeAlgorithm(context.Background(), algo.UserID).
				Return(nil, nil).
				Times(1)

			req := &desc.DescribeAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: id,
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
			const id = 1

			algo := algorithm.CreateSimpleAlgorithm(id)

			mockRepo.EXPECT().
				DescribeAlgorithm(context.Background(), algo.UserID).
				Return(nil, errors.New("some error")).
				Times(1)

			req := &desc.DescribeAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: id,
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
			const id = 0

			req := &desc.DescribeAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: id,
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
			const id = 2147483648

			req := &desc.DescribeAlgorithmRequestV1{
				Body: &desc.AlgorithmIdV1{
					Id: id,
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
				ListAlgorithms(context.Background(), uint64(limit), uint64(offset)).
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
				ListAlgorithms(context.Background(), uint64(limit), uint64(offset)).
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
				ListAlgorithms(context.Background(), uint64(limit), uint64(offset)).
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
				RemoveAlgorithm(context.Background(), uint64(id)).
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
				RemoveAlgorithm(context.Background(), uint64(id)).
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
			Expect(res).To(BeNil())
		})
	})

	When("database remove can find entity", func() {
		It("it should return nil error", func() {
			const id = 3

			notifyMock.ExpectInputWithCheckerFunctionAndSucceed(
				createAlgorithmNotificationChecker(id, notification.OpDelete))

			mockRepo.EXPECT().
				RemoveAlgorithm(context.Background(), uint64(id)).
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
			Expect(res).To(BeNil())
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
			Expect(res).To(BeNil())
		})
	})

	When("entity update got database error while request", func() {
		It("it should return Unavailable", func() {
			algo := algorithm.CreateSimpleAlgorithm(1)

			mockRepo.EXPECT().
				UpdateAlgorithm(context.Background(), algo).
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
			Expect(res).To(BeNil())
		})
	})

	When("no entity has been found to for update", func() {
		It("it should return NotFound", func() {
			algo := algorithm.CreateSimpleAlgorithm(1)

			mockRepo.EXPECT().
				UpdateAlgorithm(context.Background(), algo).
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
			Expect(res).To(BeNil())
		})
	})

	When("update were able to find entity", func() {
		It("it should return nil error", func() {
			const id = 1

			algo := algorithm.CreateSimpleAlgorithm(1)

			notifyMock.ExpectInputWithCheckerFunctionAndSucceed(
				createAlgorithmNotificationChecker(id, notification.OpUpdate))

			mockRepo.EXPECT().
				UpdateAlgorithm(context.Background(), algo).
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
			Expect(res).To(BeNil())
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
			Expect(res).To(BeNil())
		})
	})

	When("0 batchSize is provided", func() {
		It("should return OutOfRange", func() {
			algo := algorithm.CreateSimpleAlgorithm(1)

			req := &desc.MultiCreateAlgorithmRequestV1{
				Pack: []*desc.AlgorithmValueV1{
					{
						Subject:     algo.Subject,
						Description: algo.Description,
					},
				},
				BatchSize: 0,
			}

			res, err := s.MultiCreateAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.OutOfRange, "batch size (0) should be more that zero")))
			Expect(res).ToNot(BeNil())
		})
	})

	When("empty pack is provided", func() {
		It("should return InvalidArgument", func() {
			req := &desc.MultiCreateAlgorithmRequestV1{
				BatchSize: 1,
			}

			res, err := s.MultiCreateAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.InvalidArgument, "pack cannot be empty")))
			Expect(res).ToNot(BeNil())
		})
	})

	When("batch size exceeds pack size", func() {
		It("should return InvalidArgument", func() {
			algo := algorithm.CreateSimpleAlgorithm(1)

			req := &desc.MultiCreateAlgorithmRequestV1{
				Pack: []*desc.AlgorithmValueV1{
					{
						Subject:     algo.Subject,
						Description: algo.Description,
					},
				},
				BatchSize: 2,
			}

			res, err := s.MultiCreateAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.InvalidArgument, "batch size (2) should be less that size of pack(1)")))
			Expect(res).ToNot(BeNil())
		})
	})

	When("database fails all requests", func() {
		It("should return Unavailable with all request failed", func() {
			algos1_3 := createAlgorithmRangeInclusiveZeroID(1, 3)

			gomock.InOrder(
				mockRepo.EXPECT().
					AddAlgorithms(context.Background(), algos1_3[0:2]).
					Return(errors.New("some error")).
					Times(1),

				mockRepo.EXPECT().
					AddAlgorithms(context.Background(), algos1_3[2:3]).
					Return(errors.New("some more error")).
					Times(1),
			)

			req := &desc.MultiCreateAlgorithmRequestV1{
				Pack:      convertAlgorithmListToAlgorithmValueV1List(algos1_3),
				BatchSize: 2,
			}

			res, err := s.MultiCreateAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.Unavailable, "database issue: request failed")))
			Expect(res).ToNot(BeNil())
			Expect(len(res.SucceededBatches)).To(Equal(0))
		})
	})

	When("database fails the first request of two", func() {
		It("should return Unavailable with partially completed", func() {
			algos1_3 := createAlgorithmRangeInclusiveZeroID(1, 3)

			notifyMock.ExpectInputWithCheckerFunctionAndSucceed(
				createAlgorithmNotificationChecker(3, notification.OpCreate))

			gomock.InOrder(
				mockRepo.EXPECT().
					AddAlgorithms(context.Background(), algos1_3[0:2]).
					Return(errors.New("some error")).
					Times(1),

				mockRepo.EXPECT().
					AddAlgorithms(context.Background(), algos1_3[2:3]).
					DoAndReturn(func(_ context.Context, algos []algorithm.Algorithm) interface{} {
						algos[0].UserID = 3
						return nil
					}).
					Times(1),
			)

			req := &desc.MultiCreateAlgorithmRequestV1{
				Pack:      convertAlgorithmListToAlgorithmValueV1List(algos1_3),
				BatchSize: 2,
			}

			res, err := s.MultiCreateAlgorithmV1(context.Background(), req)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(status.Error(codes.Unavailable, "database issue: request partially succeeded")))
			Expect(res).ToNot(BeNil())

			succeededBunches := []*desc.AlgorithmIdPackV1{
				{
					PackIdx: 1,
					Ids: []*desc.AlgorithmIdV1{
						{
							Id: 3,
						},
					},
				},
			}
			Expect(res.SucceededBatches).To(Equal(succeededBunches))
		})
	})

	When("all database request succeeded", func() {
		It("should return Success code with all bunches", func() {
			algos1_3 := createAlgorithmRangeInclusiveZeroID(1, 3)

			// 3 times in row
			notifyMock.ExpectInputWithCheckerFunctionAndSucceed(
				createAlgorithmNotificationChecker(1, notification.OpCreate))
			notifyMock.ExpectInputWithCheckerFunctionAndSucceed(
				createAlgorithmNotificationChecker(2, notification.OpCreate))
			notifyMock.ExpectInputWithCheckerFunctionAndSucceed(
				createAlgorithmNotificationChecker(3, notification.OpCreate))

			gomock.InOrder(
				mockRepo.EXPECT().
					AddAlgorithms(context.Background(), algos1_3[0:2]).
					DoAndReturn(func(_ context.Context, algos []algorithm.Algorithm) interface{} {
						algos[0].UserID = 1
						algos[1].UserID = 2
						return nil
					}).
					Return(nil).
					Times(1),

				mockRepo.EXPECT().
					AddAlgorithms(context.Background(), algos1_3[2:3]).
					DoAndReturn(func(_ context.Context, algos []algorithm.Algorithm) interface{} {
						algos[0].UserID = 3
						return nil
					}).
					Times(1),
			)

			req := &desc.MultiCreateAlgorithmRequestV1{
				Pack:      convertAlgorithmListToAlgorithmValueV1List(algos1_3),
				BatchSize: 2,
			}

			res, err := s.MultiCreateAlgorithmV1(context.Background(), req)

			Expect(err).NotTo(HaveOccurred())
			Expect(res).ToNot(BeNil())

			succeededBunches := []*desc.AlgorithmIdPackV1{
				{
					PackIdx: 0,
					Ids: []*desc.AlgorithmIdV1{
						{
							Id: 1,
						},
						{
							Id: 2,
						},
					},
				},
				{
					PackIdx: 1,
					Ids: []*desc.AlgorithmIdV1{
						{
							Id: 3,
						},
					},
				},
			}
			Expect(res.SucceededBatches).To(Equal(succeededBunches))
		})
	})
})

func createAlgorithmRangeInclusiveZeroID(begin, end int) []algorithm.Algorithm {
	list := algorithm.CreateSimpleAlgorithmListRangeInclusive(begin, end)
	// clear ids
	for i := 0; i < len(list); i++ {
		list[i].UserID = 0
	}
	return list
}

func convertAlgorithmListToAlgorithmValueV1List(input []algorithm.Algorithm) []*desc.AlgorithmValueV1 {
	list := make([]*desc.AlgorithmValueV1, 0, len(input))
	for i := 0; i < len(input); i++ {
		list = append(list, &desc.AlgorithmValueV1{
			Subject:     input[i].Subject,
			Description: input[i].Description,
		})
	}
	return list
}

func createAlgorithmNotificationChecker(id uint64, op notification.CurOperation) func([]byte) error {
	return func(received []byte) error {
		notify := notification.NewCurNotification(id, op)
		expected, err := notify.Encode()
		if err != nil {
			return fmt.Errorf("cannot Encde() nofity: %w", err)
		}
		if bytes.Equal(expected, received) {
			return nil
		}
		return fmt.Errorf("notification are not equal to expected\n expected: %v\nreceived: %v",
			expected, received)
	}
}
