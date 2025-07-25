package v1

import (
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/shared/pkg/proto/ufo/v1"
	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/service"
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
