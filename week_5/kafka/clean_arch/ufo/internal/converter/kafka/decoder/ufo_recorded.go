package decoder

import (
	"fmt"
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"

	eventsV1 "github.com/olezhek28/microservices-course-examples/week_5/clean_arch/shared/pkg/proto/events/v1"
	"github.com/olezhek28/microservices-course-examples/week_5/clean_arch/ufo/internal/model"
)

type decoder struct{}

func NewUFORecordedDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.UFORecordedEvent, error) {
	var pb eventsV1.UFORecorded
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.UFORecordedEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	var observedAt *time.Time
	if pb.ObservedAt != nil {
		observedAt = lo.ToPtr(pb.ObservedAt.AsTime())
	}

	return model.UFORecordedEvent{
		UUID:        pb.Uuid,
		ObservedAt:  observedAt,
		Location:    pb.Location,
		Description: pb.Description,
	}, nil
}
