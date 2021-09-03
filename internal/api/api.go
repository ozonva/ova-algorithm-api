package api

import (
	"context"
	"fmt"
	"github.com/ozonva/ova-algorithm-api/internal/algorithm"
	"github.com/ozonva/ova-algorithm-api/internal/repo"
	desc "github.com/ozonva/ova-algorithm-api/pkg/ova-algorithm-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type api struct {
	desc.UnimplementedOvaAlgorithmApiServer
	repo repo.Repo
}

func (a *api) CreateAlgorithmV1(
	ctx context.Context,
	req *desc.CreateAlgorithmRequestV1,
) (*emptypb.Empty, error) {
	body := req.Body
	log.Debug().
		Str("subject", body.Subject).
		Str("description", body.Description).
		Msg("CreateAlgorithmV1 request")

	algo := algorithm.Algorithm{
		UserID:      0,
		Subject:     body.Subject,
		Description: body.Description,
	}

	if err := a.repo.AddAlgorithms([]algorithm.Algorithm{algo}); err != nil {
		log.Warn().Err(err).
			Str("subject", body.Subject).
			Str("description", body.Description).
			Msg("cannot add algorithm")

		return new(emptypb.Empty), status.Error(codes.Unavailable, "database store failed")
	}

	return new(emptypb.Empty), status.Error(codes.OK, "")
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
		return new(emptypb.Empty), status.Error(codes.OutOfRange, fmt.Sprintf("id %v", err.Error()))
	}

	found, err := a.repo.RemoveAlgorithm(id)

	if err != nil {
		log.Warn().Err(err).Msg("error occurred while RemoveAlgorithms")
		return new(emptypb.Empty), status.Error(codes.Unavailable, "database delete error")
	}

	if !found {
		return new(emptypb.Empty), status.Error(codes.NotFound, "identity not found")
	}

	return new(emptypb.Empty), status.Error(codes.OK, "successfully removed")
}

func NewOvaAlgorithmApi(repo repo.Repo) desc.OvaAlgorithmApiServer {
	return &api{repo: repo}
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