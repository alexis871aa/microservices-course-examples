package ufo

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/olezhek28/microservices-course-examples/week_4/e2e_tests/ufo/internal/repository"
)

var _ def.UFORepository = (*repository)(nil)

const (
	collectionName = "sightings"
)

type repository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client, databaseName string) *repository {
	return &repository{
		client:     client,
		collection: client.Database(databaseName).Collection(collectionName),
	}
}
