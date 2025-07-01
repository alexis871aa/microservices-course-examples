package converter

import (
	"context"
	"time"

	"github.com/samber/lo"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/model"
	repoModel "github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/repository/model"
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

// SightingToRedisView - конвертер из модели домена в Redis view
func SightingToRedisView(ctx context.Context, sighting model.Sighting) repoModel.SightingRedisView {
	var observedAt *int64
	if sighting.Info.ObservedAt != nil {
		observedAt = lo.ToPtr(sighting.Info.ObservedAt.UnixNano())
	}

	var updatedAt *int64
	if sighting.UpdatedAt != nil {
		updatedAt = lo.ToPtr(sighting.UpdatedAt.UnixNano())
	}

	var deletedAt *int64
	if sighting.DeletedAt != nil {
		deletedAt = lo.ToPtr(sighting.DeletedAt.UnixNano())
	}

	return repoModel.SightingRedisView{
		UUID:         sighting.Uuid,
		ObservedAtNs: observedAt,
		Location:     sighting.Info.Location,
		Description:  sighting.Info.Description,
		Color:        sighting.Info.Color,
		Sound:        sighting.Info.Sound,
		Duration:     sighting.Info.DurationSeconds,
		CreatedAtNs:  sighting.CreatedAt.UnixNano(),
		UpdatedAtNs:  updatedAt,
		DeletedAtNs:  deletedAt,
	}
}

// SightingFromRedisView - конвертер из Redis view в модель домена
func SightingFromRedisView(ctx context.Context, redisView repoModel.SightingRedisView) model.Sighting {
	var observedAt *time.Time
	if redisView.ObservedAtNs != nil {
		tmp := time.Unix(0, *redisView.ObservedAtNs)
		observedAt = &tmp
	}

	var updatedAt *time.Time
	if redisView.UpdatedAtNs != nil {
		tmp := time.Unix(0, *redisView.UpdatedAtNs)
		updatedAt = &tmp
	}

	var deletedAt *time.Time
	if redisView.DeletedAtNs != nil {
		tmp := time.Unix(0, *redisView.DeletedAtNs)
		deletedAt = &tmp
	}

	return model.Sighting{
		Uuid: redisView.UUID,
		Info: model.SightingInfo{
			ObservedAt:      observedAt,
			Location:        redisView.Location,
			Description:     redisView.Description,
			Color:           redisView.Color,
			Sound:           redisView.Sound,
			DurationSeconds: redisView.Duration,
		},
		CreatedAt: time.Unix(0, redisView.CreatedAtNs),
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}
