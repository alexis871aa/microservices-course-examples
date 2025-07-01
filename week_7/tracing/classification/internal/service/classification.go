package service

import (
	"context"
	"strings"

	"github.com/olezhek28/microservices-course-examples/week_7/tracing/platform/pkg/tracing"
)

// ClassificationService - сервис для классификации НЛО
type ClassificationService struct{}

// NewClassificationService создает новый сервис классификации
func NewClassificationService() *ClassificationService {
	return &ClassificationService{}
}

// ClassifyObject классифицирует объект на основе описания, цвета и длительности
func (s *ClassificationService) ClassifyObject(ctx context.Context, description, color string, durationSeconds int32) (string, float32, string, error) {
	// Создаем спан для внутренней логики
	ctx, span := tracing.StartSpan(ctx, "classification.analyze")
	defer span.End()

	// Простая логика классификации
	objectType := "unknown"
	confidence := float32(0.5)
	explanation := "Базовая классификация"

	desc := strings.ToLower(description)
	color = strings.ToLower(color)

	// Классификация по форме
	if strings.Contains(desc, "треугольник") || strings.Contains(desc, "triangle") {
		objectType = "triangular_craft"
		confidence = 0.8
		explanation = "Треугольная форма характерна для современных НЛО"
	} else if strings.Contains(desc, "диск") || strings.Contains(desc, "disk") || strings.Contains(desc, "тарелка") {
		objectType = "classic_saucer"
		confidence = 0.9
		explanation = "Классическая форма летающей тарелки"
	} else if strings.Contains(desc, "сфера") || strings.Contains(desc, "шар") || strings.Contains(desc, "sphere") {
		objectType = "orb"
		confidence = 0.7
		explanation = "Сферический объект неизвестного происхождения"
	}

	// Корректировка по цвету
	switch color {
	case "зеленый", "green":
		confidence += 0.1
		explanation += ". Зеленое свечение часто наблюдается у НЛО"
	case "красный", "red":
		confidence += 0.05
		explanation += ". Красный цвет может указывать на двигательную систему"
	case "белый", "white":
		confidence += 0.02
		explanation += ". Белое свечение - распространенное явление"
	}

	// Корректировка по длительности
	if durationSeconds > 600 { // более 10 минут
		confidence += 0.1
		explanation += ". Длительное наблюдение повышает достоверность"
	} else if durationSeconds < 30 { // менее 30 секунд
		confidence -= 0.1
		explanation += ". Кратковременное наблюдение снижает достоверность"
	}

	// Ограничиваем confidence в пределах [0, 1]
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	return objectType, confidence, explanation, nil
}
