package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	ufoV1 "github.com/olezhek28/microservices-course-examples/week_4/e2e_tests/shared/pkg/proto/ufo/v1"
	"github.com/olezhek28/microservices-course-examples/week_4/e2e_tests/ufo/internal/model"
)

func (a *api) Delete(ctx context.Context, req *ufoV1.DeleteRequest) (*emptypb.Empty, error) {
	err := a.ufoService.Delete(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrSightingNotFound) {
			return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
