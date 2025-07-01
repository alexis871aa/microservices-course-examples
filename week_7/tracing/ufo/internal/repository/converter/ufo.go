package converter

import (
	"github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/model"
	repoModel "github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/repository/model"
)

func SightingInfoToRepoModel(info model.SightingInfo) repoModel.SightingInfo {
	return repoModel.SightingInfo{
		ObservedAt:      info.ObservedAt,
		Location:        info.Location,
		Description:     info.Description,
		Color:           info.Color,
		Sound:           info.Sound,
		DurationSeconds: info.DurationSeconds,
	}
}

func SightingToModel(sighting repoModel.Sighting) model.Sighting {
	return model.Sighting{
		Uuid:      sighting.Uuid,
		Info:      SightingInfoToModel(sighting.Info),
		CreatedAt: sighting.CreatedAt,
		UpdatedAt: sighting.UpdatedAt,
		DeletedAt: sighting.DeletedAt,
	}
}

func SightingInfoToModel(info repoModel.SightingInfo) model.SightingInfo {
	return model.SightingInfo{
		ObservedAt:      info.ObservedAt,
		Location:        info.Location,
		Description:     info.Description,
		Color:           info.Color,
		Sound:           info.Sound,
		DurationSeconds: info.DurationSeconds,
	}
}
