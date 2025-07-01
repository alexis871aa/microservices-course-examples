package kafka

import "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/model"

type UFORecordedDecoder interface {
	Decode(data []byte) (model.UFORecordedEvent, error)
}
