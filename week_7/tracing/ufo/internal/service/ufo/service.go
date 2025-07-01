package ufo

import (
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/client/grpc"
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/repository"
	def "github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/service"
)

var _ def.UFOService = (*service)(nil)

type service struct {
	ufoRepository  repository.UFORepository
	analysisClient grpc.AnalysisClient
}

func NewService(ufoRepository repository.UFORepository, analysisClient grpc.AnalysisClient) *service {
	return &service{
		ufoRepository:  ufoRepository,
		analysisClient: analysisClient,
	}
}
