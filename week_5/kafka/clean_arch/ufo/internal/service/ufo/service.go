package ufo

import (
	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/repository"
	def "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/service"
)

var _ def.UFOService = (*service)(nil)

type service struct {
	ufoRepository      repository.UFORepository
	ufoProducerService def.UFOProducerService
}

func NewService(ufoRepository repository.UFORepository, ufoProducerService def.UFOProducerService) *service {
	return &service{
		ufoRepository:      ufoRepository,
		ufoProducerService: ufoProducerService,
	}
}
