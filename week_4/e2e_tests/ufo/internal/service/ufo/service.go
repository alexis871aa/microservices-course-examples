package ufo

import (
	"github.com/olezhek28/microservices-course-examples/week_4/e2e_tests/ufo/internal/repository"
	def "github.com/olezhek28/microservices-course-examples/week_4/e2e_tests/ufo/internal/service"
)

var _ def.UFOService = (*service)(nil)

type service struct {
	ufoRepository repository.UFORepository
}

func NewService(ufoRepository repository.UFORepository) *service {
	return &service{
		ufoRepository: ufoRepository,
	}
}
