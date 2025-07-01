package ufo

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/olezhek28/microservices-course-examples/week_5/telegram/clean_arch/ufo/internal/repository"
)

var _ def.UFORepository = (*repository)(nil)

const (
	collectionName = "sightings"
)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	repo := &repository{
		collection: db.Collection(collectionName),
	}

	return repo
}
