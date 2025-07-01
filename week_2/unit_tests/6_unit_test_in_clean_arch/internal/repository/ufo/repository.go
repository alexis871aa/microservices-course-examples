package ufo

import (
	"sync"

	def "github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/repository"
	repoModel "github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/repository/model"
)

var _ def.UFORepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Sighting
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]repoModel.Sighting),
	}
}
