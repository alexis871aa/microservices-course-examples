package v1

import (
	"github.com/olezhek28/microservices-course-examples/week_2/layers/internal/service"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_2/layers/pkg/proto/ufo/v1"
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
