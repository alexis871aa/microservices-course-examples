package ufo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/model"
)

func (r *repository) Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error {
	// Проверяем существование документа
	var existing bson.M
	err := r.collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&existing)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.ErrSightingNotFound
		}
		return err
	}

	// Формируем update запрос
	updateDoc := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Обновляем поля, только если они были установлены в запросе
	if updateInfo.ObservedAt != nil {
		updateDoc["$set"].(bson.M)["info.observed_at"] = updateInfo.ObservedAt
	}

	if updateInfo.Location != nil {
		updateDoc["$set"].(bson.M)["info.location"] = *updateInfo.Location
	}

	if updateInfo.Description != nil {
		updateDoc["$set"].(bson.M)["info.description"] = *updateInfo.Description
	}

	if updateInfo.Color != nil {
		updateDoc["$set"].(bson.M)["info.color"] = updateInfo.Color
	}

	if updateInfo.Sound != nil {
		updateDoc["$set"].(bson.M)["info.sound"] = updateInfo.Sound
	}

	if updateInfo.DurationSeconds != nil {
		updateDoc["$set"].(bson.M)["info.duration_seconds"] = updateInfo.DurationSeconds
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": uuid}, updateDoc)
	if err != nil {
		return err
	}

	return nil
}
