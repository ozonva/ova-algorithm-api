package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"github.com/ozonva/ova-algorithm-api/internal/notification"
	"github.com/ozonva/ova-algorithm-api/internal/repo"
	desc "github.com/ozonva/ova-algorithm-api/pkg/ova-algorithm-api"
)

var regCounterCreateAlgorithm = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "algorithm_created_notifications",
	Help: "Notifications send to Kafka for Algorithms creations",
})

var regCounterUpdateAlgorithm = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "algorithm_updated_notifications",
	Help: "Notifications send to Kafka for Algorithms updates",
})

var regCounterDeleteAlgorithm = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "algorithm_deleted_notifications",
	Help: "Notifications send to Kafka for Algorithms deletes",
})

func init() {
	prometheus.MustRegister(
		regCounterCreateAlgorithm,
		regCounterUpdateAlgorithm,
		regCounterDeleteAlgorithm,
	)
}

const defaultTopicName = "ova.algorithm.notify"

type api struct {
	desc.UnimplementedOvaAlgorithmApiServer
	repo     repo.Repo
	producer sarama.AsyncProducer
}

func (a *api) CreateAlgorithmV1(
	ctx context.Context,
	req *desc.CreateAlgorithmRequestV1,
) (*desc.CreateAlgorithmResponseV1, error) {
	body := req.Body
	log.Debug().
		Str("subject", body.Subject).
		Str("description", body.Description).
		Msg("CreateAlgorithmV1 request")

	algos := []algorithm.Algorithm{
		{
			UserID:      0,
			Subject:     body.Subject,
			Description: body.Description,
		},
	}

	res := &desc.CreateAlgorithmResponseV1{}

	if err := a.repo.AddAlgorithms(algos); err != nil {
		log.Warn().Err(err).
			Str("subject", body.Subject).
			Str("description", body.Description).
			Msg("cannot add algorithm")

		return res, status.Error(codes.Unavailable, "database store failed")
	}

	res.Body = &desc.AlgorithmIdV1{
		Id: int64(algos[0].UserID),
	}

	log.Debug().
		Uint64("UserID", algos[0].UserID).
		Msg("CreateAlgorithmV1 UserID assigned")

	a.notifyKafkaAlgorithmPack(algos, notification.OpCreate)
	regCounterCreateAlgorithm.Inc()

	return res, status.Error(codes.OK, "")
}

func (a *api) DescribeAlgorithmV1(
	ctx context.Context,
	req *desc.DescribeAlgorithmRequestV1,
) (*desc.DescribeAlgorithmResponseV1, error) {
	log.Debug().
		Int64("id", req.Body.Id).
		Msg("DescribeAlgorithmV1 request")

	id, err := validateOneInt32MaxRangeInt64(req.Body.Id)
	if err != nil {
		return nil, status.Error(codes.OutOfRange, fmt.Sprintf("id %v", err.Error()))
	}

	algo, err := a.repo.DescribeAlgorithm(id)

	if err != nil {
		return nil, status.Error(codes.Unavailable, "database fetch error")
	}

	if algo == nil {
		return nil, status.Error(codes.NotFound, "identity not found")
	}

	res := &desc.DescribeAlgorithmResponseV1{
		Body: &desc.AlgorithmV1{
			Id:          int64(algo.UserID),
			Subject:     algo.Subject,
			Description: algo.Description,
		},
	}
	return res, status.Error(codes.OK, "")
}

func (a *api) ListAlgorithmsV1(
	ctx context.Context,
	req *desc.ListAlgorithmsRequestV1,
) (*desc.ListAlgorithmsResponseV1, error) {
	log.Debug().
		Int64("offset", req.Offset.Id).
		Int64("limit", req.Limit).
		Msg("ListAlgorithmsV1 request")

	id, err := validateZeroInt32MaxRangeInt64(req.Offset.Id)
	if err != nil {
		return nil, status.Error(codes.OutOfRange, fmt.Sprintf("offset %v", err.Error()))
	}

	limit, err := validateOneInt32MaxRangeInt64(req.Limit)
	if err != nil {
		return nil, status.Error(codes.OutOfRange, fmt.Sprintf("limit %v", err.Error()))
	}

	list, err := a.repo.ListAlgorithms(limit, id)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "database list error")
	}

	if len(list) == 0 {
		return nil, status.Error(codes.NotFound, "entities not found")
	}

	idList := make([]*desc.AlgorithmV1, 0, len(list))
	for i := 0; i < len(list); i++ {
		idList = append(idList, &desc.AlgorithmV1{
			Id:          int64(list[i].UserID),
			Subject:     list[i].Subject,
			Description: list[i].Description,
		})
	}

	res := &desc.ListAlgorithmsResponseV1{
		Body: idList,
	}

	return res, status.Error(codes.OK, "")
}

func (a *api) RemoveAlgorithmV1(
	ctx context.Context,
	req *desc.RemoveAlgorithmRequestV1,
) (*emptypb.Empty, error) {
	log.Debug().
		Int64("id", req.Body.Id).
		Msg("RemoveAlgorithmV1")

	id, err := validateOneInt32MaxRangeInt64(req.Body.Id)
	if err != nil {
		return nil, status.Error(codes.OutOfRange, fmt.Sprintf("id %v", err.Error()))
	}

	found, err := a.repo.RemoveAlgorithm(id)

	if err != nil {
		log.Warn().Err(err).Msg("error occurred while RemoveAlgorithms")
		return nil, status.Error(codes.Unavailable, "database delete error")
	}

	if !found {
		return nil, status.Error(codes.NotFound, "identity not found")
	}

	a.notifyKafkaAlgorithmOne(id, notification.OpDelete)
	regCounterDeleteAlgorithm.Inc()

	return new(emptypb.Empty), status.Error(codes.OK, "successfully removed")
}

func (a *api) UpdateAlgorithmV1(
	ctx context.Context,
	req *desc.UpdateAlgorithmRequestV1,
) (*emptypb.Empty, error) {
	log.Debug().
		Int64("id", req.Body.Id).
		Msg("UpdateAlgorithmV1")

	id, err := validateOneInt32MaxRangeInt64(req.Body.Id)
	if err != nil {
		return nil, status.Error(codes.OutOfRange, fmt.Sprintf("id %v", err.Error()))
	}

	entity := algorithm.Algorithm{
		UserID:      id,
		Subject:     req.Body.Subject,
		Description: req.Body.Description,
	}

	found, err := a.repo.UpdateAlgorithm(entity)

	if err != nil {
		log.Warn().Err(err).Msg("error occurred while UpdateAlgorithmV1")
		return nil, status.Error(codes.Unavailable, "database update error")
	}

	if !found {
		return nil, status.Error(codes.NotFound, "identity not found")
	}

	a.notifyKafkaAlgorithmOne(entity.UserID, notification.OpUpdate)
	regCounterUpdateAlgorithm.Inc()

	return new(emptypb.Empty), status.Error(codes.OK, "successfully updated")
}

func (a *api) MultiCreateAlgorithmV1(
	ctx context.Context,
	req *desc.MultiCreateAlgorithmRequestV1,
) (*desc.MultiCreateAlgorithmResponseV1, error) {
	parentSpan, _ := opentracing.StartSpanFromContext(ctx, "MultiCreateAlgorithmV1")
	defer parentSpan.Finish()

	parentSpan.LogKV("batchSize", req.BatchSize)

	log.Debug().
		Int32("batchSize", req.BatchSize).
		Int("len(pack)", len(req.Pack)).
		Msg("MultiCreateAlgorithmV1")

	if req.BatchSize < 1 {
		return new(desc.MultiCreateAlgorithmResponseV1), status.Errorf(
			codes.OutOfRange,
			fmt.Sprintf("batch size (%v) should be more that zero", req.BatchSize),
		)
	}

	packSize := len(req.Pack)
	if packSize == 0 {
		return new(desc.MultiCreateAlgorithmResponseV1), status.Errorf(
			codes.InvalidArgument,
			"pack cannot be empty",
		)
	}

	if int(req.BatchSize) > packSize {
		return new(desc.MultiCreateAlgorithmResponseV1), status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("batch size (%v) should be less that size of pack(%v)", req.BatchSize, packSize),
		)
	}

	algos := make([]algorithm.Algorithm, 0, packSize)
	for i := 0; i < packSize; i++ {
		algos = append(algos, algorithm.Algorithm{
			UserID:      0,
			Subject:     req.Pack[i].Subject,
			Description: req.Pack[i].Description,
		})
	}

	algoPacks := algorithm.SplitAlgorithmsToBulks(algos, uint(req.BatchSize))
	succeededBatches := make([]*desc.AlgorithmIdPackV1, 0, len(algoPacks))

	for i := 0; i < len(algoPacks); i++ {
		func() {
			childSpan := opentracing.StartSpan(
				"AddAlgorithms",
				opentracing.ChildOf(parentSpan.Context()))

			childSpan.LogKV("batchSize", len(algoPacks[i]))

			if err := a.repo.AddAlgorithms(algoPacks[i]); err != nil {
				log.Warn().Err(err).
					Int("index", i).
					Int("batchSize", len(algoPacks[i])).
					Msg("failed to add batch")

			} else {
				succeededBatches = append(succeededBatches, createAlgorithmIDPackV1(i, algoPacks[i]))
				a.notifyKafkaAlgorithmPack(algoPacks[i], notification.OpCreate)
				regCounterCreateAlgorithm.Add(float64(len(algoPacks[i])))
			}

			defer childSpan.Finish()
		}()
	}

	res := &desc.MultiCreateAlgorithmResponseV1{
		SucceededBatches: succeededBatches,
	}

	if len(succeededBatches) == 0 {
		return res, status.Error(codes.Unavailable, "database issue: request failed")
	}

	if len(succeededBatches) != len(algoPacks) {
		return res, status.Error(codes.Unavailable, "database issue: request partially succeeded")
	}

	return res, status.Error(codes.OK, "")
}

func createAlgorithmIDPackV1(packIdx int, list []algorithm.Algorithm) *desc.AlgorithmIdPackV1 {
	ids := make([]*desc.AlgorithmIdV1, 0, len(list))
	for i := 0; i < len(list); i++ {
		ids = append(ids, &desc.AlgorithmIdV1{
			Id: int64(list[i].UserID),
		})
	}
	return &desc.AlgorithmIdPackV1{
		PackIdx: int32(packIdx),
		Ids:     ids,
	}
}

func NewOvaAlgorithmAPI(repo repo.Repo, producer sarama.AsyncProducer) desc.OvaAlgorithmApiServer {
	return &api{
		repo:     repo,
		producer: producer,
	}
}

func (a *api) notifyKafkaAlgorithmOne(id uint64, op notification.CurOperation) {
	opNotification := notification.NewCurNotification(id, op)

	a.producer.Input() <- &sarama.ProducerMessage{
		Topic: defaultTopicName,
		Key:   sarama.ByteEncoder([]byte(strconv.FormatUint(id, 10))),
		Value: opNotification,
	}
}

func (a *api) notifyKafkaAlgorithmPack(list []algorithm.Algorithm, op notification.CurOperation) {
	for i := 0; i < len(list); i++ {
		a.notifyKafkaAlgorithmOne(list[i].UserID, op)
	}
}

func validateOneInt32MaxRangeInt64(id int64) (uint64, error) {
	//postgres SERIAL range 1 - 2,147,483,647
	if id < 1 || id > 2147483647 {
		return 0, fmt.Errorf("(%v) is out of range 1 - 2,147,483,647", id)
	}
	return uint64(id), nil
}

func validateZeroInt32MaxRangeInt64(id int64) (uint64, error) {
	//postgres LIMIT range 0 - 2,147,483,647
	if id < 0 || id > 2147483647 {
		return 0, fmt.Errorf("(%v) is out of range 0 - 2,147,483,647", id)
	}
	return uint64(id), nil
}
