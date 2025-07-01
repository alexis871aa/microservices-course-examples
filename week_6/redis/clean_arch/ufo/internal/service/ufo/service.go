package ufo

import (
	"time"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/repository"
	def "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/service"
)

var _ def.UFOService = (*service)(nil)

type service struct {
	ufoRepository   repository.UFORepository
	cacheRepository repository.UFOCacheRepository
	cacheTTL        time.Duration
}

func NewService(
	ufoRepository repository.UFORepository,
	cacheRepository repository.UFOCacheRepository,
	cacheTTL time.Duration,
) *service {
	return &service{
		ufoRepository:   ufoRepository,
		cacheRepository: cacheRepository,
		cacheTTL:        cacheTTL,
	}
}
