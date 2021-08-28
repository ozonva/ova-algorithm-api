package api

import (
	"context"
	desc "github.com/ozonva/ova-algorithm-api/pkg/ova-algorithm-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	errNotImplemented = "method is not implemented"
)

type api struct {
	desc.UnimplementedOvaAlgorithmApiServer
}

func (a *api) CreateAlgorithmV1(
	ctx context.Context,
	req *desc.CreateAlgorithmRequestV1,
) (*desc.CreateAlgorithmResponseV1, error) {
	body := req.Body
	log.Info().
		Str("subject", body.Subject).
		Str("description", body.Description).
		Msg("CreateAlgorithmV1")

	return nil, status.Error(codes.Unimplemented, errNotImplemented)
}

func (a *api) DescribeAlgorithmV1(
	ctx context.Context,
	req *desc.DescribeAlgorithmRequestV1,
	) (*desc.DescribeAlgorithmResponseV1, error) {
	log.Info().
		Int64("id", req.Body.Id).
		Msg("DescribeAlgorithmV1")
	return nil, status.Error(codes.Unimplemented, errNotImplemented)
}

func (a *api) ListAlgorithmsV1(
	ctx context.Context,
	req *desc.ListAlgorithmsRequestV1,
) (*desc.ListAlgorithmsResponseV1, error) {
	log.Info().
		Int64("offset", req.Offset.Id).
		Int64("limit", req.Limit).
		Msg("ListAlgorithmsV1")
	return nil, status.Error(codes.NotFound, errNotImplemented)
}

func (a *api) RemoveAlgorithmV1(
	ctx context.Context,
	req *desc.RemoveAlgorithmRequestV1,
) (*emptypb.Empty, error) {
	log.Info().
		Int64("id", req.Body.Id).
		Msg("RemoveAlgorithmV1")
	return nil, status.Error(codes.NotFound, errNotImplemented)
}

func NewOvaAlgorithmApi() desc.OvaAlgorithmApiServer {
	return &api{
	}
}
