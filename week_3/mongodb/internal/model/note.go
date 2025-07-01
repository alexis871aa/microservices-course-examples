package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Note - модель заметки в MongoDB
type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Body      string             `bson:"body"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at,omitempty"`
}
