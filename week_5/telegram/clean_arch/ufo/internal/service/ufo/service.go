package ufo

import (
	"github.com/olezhek28/microservices-course-examples/week_5/telegram/clean_arch/ufo/internal/repository"
	def "github.com/olezhek28/microservices-course-examples/week_5/telegram/clean_arch/ufo/internal/service"
)

var _ def.UFOService = (*service)(nil)

type service struct {
	ufoRepository   repository.UFORepository
	telegramService def.TelegramService
}

func NewService(ufoRepository repository.UFORepository, telegramService def.TelegramService) *service {
	return &service{
		ufoRepository:   ufoRepository,
		telegramService: telegramService,
	}
}
