package v1

import (
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_4/di/shared/pkg/proto/ufo/v1"
	"github.com/olezhek28/microservices-course-examples/week_4/di/ufo/internal/service"
)

type api struct {
	ufoV1.UnimplementedUFOServiceServer

	ufoService service.UFOService
}

func NewAPI(ufoService service.UFOService) *api {
	return &api{
		ufoService: ufoService,
	}
}
