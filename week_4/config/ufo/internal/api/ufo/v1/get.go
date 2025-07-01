package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ufoV1 "github.com/olezhek28/microservices-course-examples/week_4/config/shared/pkg/proto/ufo/v1"
	"github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/converter"
	"github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/model"
)

func (a *api) Get(ctx context.Context, req *ufoV1.GetRequest) (*ufoV1.GetResponse, error) {
	sighting, err := a.ufoService.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrSightingNotFound) {
			return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}

	return &ufoV1.GetResponse{
		Sighting: converter.SightingToProto(sighting),
	}, nil
}
