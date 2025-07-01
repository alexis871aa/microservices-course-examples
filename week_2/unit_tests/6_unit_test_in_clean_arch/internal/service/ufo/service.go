package ufo

import (
	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/repository"
	def "github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/service"
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
