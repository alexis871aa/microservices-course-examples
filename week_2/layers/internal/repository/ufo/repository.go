package ufo

import (
	"sync"

	def "github.com/olezhek28/microservices-course-examples/week_2/layers/internal/repository"
	repoModel "github.com/olezhek28/microservices-course-examples/week_2/layers/internal/repository/model"
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
