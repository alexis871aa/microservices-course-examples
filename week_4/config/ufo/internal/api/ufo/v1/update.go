package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	ufoV1 "github.com/olezhek28/microservices-course-examples/week_4/config/shared/pkg/proto/ufo/v1"
	"github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/converter"
	"github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/model"
)

func (a *api) Update(ctx context.Context, req *ufoV1.UpdateRequest) (*emptypb.Empty, error) {
	if req.UpdateInfo == nil {
		return nil, status.Error(codes.InvalidArgument, "update_info cannot be nil")
	}

	err := a.ufoService.Update(ctx, req.GetUuid(), converter.UpdateInfoToModel(req.GetUpdateInfo()))
	if err != nil {
		if errors.Is(err, model.ErrSightingNotFound) {
			return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
