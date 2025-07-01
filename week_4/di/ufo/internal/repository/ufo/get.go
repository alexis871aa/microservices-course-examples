package ufo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/olezhek28/microservices-course-examples/week_4/di/ufo/internal/model"
	repoConverter "github.com/olezhek28/microservices-course-examples/week_4/di/ufo/internal/repository/converter"
	repoModel "github.com/olezhek28/microservices-course-examples/week_4/di/ufo/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	var repoSighting repoModel.Sighting

	err := r.collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&repoSighting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Sighting{}, model.ErrSightingNotFound
		}
		return model.Sighting{}, err
	}

	return repoConverter.SightingToModel(repoSighting), nil
}
