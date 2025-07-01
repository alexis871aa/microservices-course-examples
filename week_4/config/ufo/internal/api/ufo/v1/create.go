package v1

import (
	"context"

	ufoV1 "github.com/olezhek28/microservices-course-examples/week_4/config/shared/pkg/proto/ufo/v1"
	"github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/converter"
)

func (a *api) Create(ctx context.Context, req *ufoV1.CreateRequest) (*ufoV1.CreateResponse, error) {
	uuid, err := a.ufoService.Create(ctx, converter.UFOInfoToModel(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &ufoV1.CreateResponse{
		Uuid: uuid,
	}, nil
}
