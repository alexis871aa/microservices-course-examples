package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/olezhek28/microservices-course-examples/week_3/mongodb/internal/model"
	"github.com/olezhek28/microservices-course-examples/week_3/mongodb/internal/repository"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	// Получаем строку подключения из переменной окружения
	dbURI := os.Getenv("MONGO_URI")

	// Создаем клиент MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		cerr := client.Disconnect(ctx)
		if cerr != nil {
			log.Printf("failed to disconnect: %v\n", cerr)
		}
	}()

	// Проверяем соединение с базой данных
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	// Получаем базу данных
	db := client.Database("example")

	// Создаем репозиторий заметок
	noteRepo := repository.NewNoteRepository(db)

	// Создаем новую заметку
	note := model.Note{
		Title:     gofakeit.City(),
		Body:      gofakeit.Address().Street,
		CreatedAt: time.Now(),
	}

	// Вставляем заметку
	insertedID, err := noteRepo.Create(ctx, note)
	if err != nil {
		log.Printf("failed to create note: %v\n", err)
		return
	}
	log.Printf("inserted note with id: %s\n", insertedID.Hex())

	// Получаем все заметки
	notes, err := noteRepo.GetAll(ctx)
	if err != nil {
		log.Printf("failed to get notes: %v\n", err)
		return
	}

	for _, n := range notes {
		updatedAtStr := "nil"
		if n.UpdatedAt != nil {
			updatedAtStr = n.UpdatedAt.String()
		}
		log.Printf("id: %s, title: %s, body: %s, created_at: %v, updated_at: %v\n",
			n.ID.Hex(), n.Title, n.Body, n.CreatedAt, updatedAtStr)
	}

	// Получаем заметку по ID
	foundNote, err := noteRepo.GetByID(ctx, insertedID)
	if err != nil {
		log.Printf("failed to get note by id: %v\n", err)
		return
	}
	log.Printf("found note - id: %s, title: %s, body: %s\n",
		foundNote.ID.Hex(), foundNote.Title, foundNote.Body)

	// Обновляем заметку
	foundNote.Title = gofakeit.City()
	foundNote.Body = gofakeit.Address().Street

	err = noteRepo.Update(ctx, foundNote)
	if err != nil {
		log.Printf("failed to update note: %v\n", err)
		return
	}
	log.Printf("updated note with id: %s\n", foundNote.ID.Hex())

	// Получаем обновленную заметку
	updatedNote, err := noteRepo.GetByID(ctx, insertedID)
	if err != nil {
		log.Printf("failed to get updated note: %v\n", err)
		return
	}
	updatedAtStr := "nil"
	if updatedNote.UpdatedAt != nil {
		updatedAtStr = updatedNote.UpdatedAt.String()
	}
	log.Printf("Updated note - id: %s, title: %s, body: %s, created_at: %v, updated_at: %v\n",
		updatedNote.ID.Hex(), updatedNote.Title, updatedNote.Body, updatedNote.CreatedAt, updatedAtStr)

	// Создаем заметку для демонстрации удаления
	noteForDelete := model.Note{
		Title:     gofakeit.City(),
		Body:      gofakeit.Address().Street,
		CreatedAt: time.Now(),
	}

	deleteID, err := noteRepo.Create(ctx, noteForDelete)
	if err != nil {
		log.Printf("failed to create note for deletion: %v\n", err)
		return
	}
	log.Printf("inserted note for deletion with id: %s\n", deleteID.Hex())

	// Удаляем заметку
	deletedCount, err := noteRepo.Delete(ctx, deleteID)
	if err != nil {
		log.Printf("failed to delete note: %v\n", err)
		return
	}
	log.Printf("deleted %d note(s) with id: %s\n", deletedCount, deleteID.Hex())
}
