package ufo

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/repository"
)

var _ def.UFORepository = (*repository)(nil)

const (
	databaseName   = "ufo_db"
	collectionName = "sightings"
)

type repository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client) *repository {
	return &repository{
		client:     client,
		collection: client.Database(databaseName).Collection(collectionName),
	}
}
