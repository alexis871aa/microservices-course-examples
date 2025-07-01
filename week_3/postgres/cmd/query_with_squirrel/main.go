package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv("DB_URI")

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer pool.Close()

	// Делаем запрос на вставку записи в таблицу note
	builderInsert := sq.Insert("note").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "body").
		Values(gofakeit.City(), gofakeit.Address().Street).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return
	}

	var noteID int
	err = pool.QueryRow(ctx, query, args...).Scan(&noteID)
	if err != nil {
		log.Printf("failed to insert note: %v\n", err)
		return
	}

	log.Printf("inserted note with id: %d\n", noteID)

	// Делаем запрос на выборку записей из таблицы note
	builderSelect := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("failed to select notes: %v\n", err)
		return
	}

	var id int
	var title, body string
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&id, &title, &body, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("failed to scan note: %v\n", err)
			return
		}

		log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updatedAt)
	}

	// Делаем запрос на обновление записи в таблице note
	builderUpdate := sq.Update("note").
		PlaceholderFormat(sq.Dollar).
		Set("title", gofakeit.City()).
		Set("body", gofakeit.Address().Street).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": noteID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update note: %v\n", err)
		return
	}

	log.Printf("updated %d rows", res.RowsAffected())

	// Делаем запрос на получение измененной записи из таблицы note
	builderSelectOne := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": noteID}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &title, &body, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("failed to select notes: %v\n", err)
		return
	}

	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updatedAt)

	// Создаем еще одну запись для демонстрации удаления
	builderInsertForDelete := sq.Insert("note").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "body").
		Values(gofakeit.City(), gofakeit.Address().Street).
		Suffix("RETURNING id")

	query, args, err = builderInsertForDelete.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return
	}

	var deleteNoteID int
	err = pool.QueryRow(ctx, query, args...).Scan(&deleteNoteID)
	if err != nil {
		log.Printf("failed to insert note for deletion: %v\n", err)
		return
	}

	log.Printf("inserted note for deletion with id: %d\n", deleteNoteID)

	// Делаем запрос на удаление записи из таблицы note
	builderDelete := sq.Delete("note").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": deleteNoteID})

	query, args, err = builderDelete.ToSql()
	if err != nil {
		log.Printf("failed to build delete query: %v\n", err)
		return
	}

	res, err = pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete note: %v\n", err)
		return
	}

	log.Printf("deleted %d rows with id: %d", res.RowsAffected(), deleteNoteID)
}
