package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/olezhek28/microservices-course-examples/week_3/mongodb/internal/model"
)

// NoteRepository представляет репозиторий для работы с заметками
type NoteRepository struct {
	collection *mongo.Collection
}

// NewNoteRepository создает новый репозиторий заметок
func NewNoteRepository(db *mongo.Database) *NoteRepository {
	collection := db.Collection("notes")

	// Создаем индексы при инициализации
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "title", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	return &NoteRepository{
		collection: collection,
	}
}

// Create создает новую заметку
func (r *NoteRepository) Create(ctx context.Context, note model.Note) (primitive.ObjectID, error) {
	if note.CreatedAt.IsZero() {
		note.CreatedAt = time.Now()
	}

	res, err := r.collection.InsertOne(ctx, note)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

// GetAll получает все заметки
func (r *NoteRepository) GetAll(ctx context.Context) ([]model.Note, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var notes []model.Note
	err = cursor.All(ctx, &notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

// GetByID получает заметку по ID
func (r *NoteRepository) GetByID(ctx context.Context, id primitive.ObjectID) (model.Note, error) {
	var note model.Note
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&note)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Note{}, model.ErrNoteNotFound
		}

		return model.Note{}, err
	}

	return note, nil
}

// Update обновляет заметку
func (r *NoteRepository) Update(ctx context.Context, note model.Note) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": note.ID},
		bson.M{
			"$set": bson.M{
				"title":      note.Title,
				"body":       note.Body,
				"updated_at": lo.ToPtr(time.Now()),
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete удаляет заметку по ID
func (r *NoteRepository) Delete(ctx context.Context, id primitive.ObjectID) (int64, error) {
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}
