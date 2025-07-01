package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/olezhek28/microservices-course-examples/week_3/postgres/internal/migrator"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv("DB_URI")

	// Создаем соединение с базой данных
	con, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		cerr := con.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close connection: %v\n", cerr)
		}
	}()

	// Проверяем, что соединение с базой установлено
	err = con.Ping(ctx)
	if err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		return
	}

	// Инициализируем мигратор
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*con.Config().Copy()), migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Ошибка миграции базы данных: %v\n", err)
		return
	}

	// Делаем запрос на вставку записи в таблицу note
	res, err := con.Exec(ctx, "INSERT INTO note (title, body) VALUES ($1, $2)", gofakeit.City(), gofakeit.Address().Street)
	if err != nil {
		log.Printf("failed to insert note: %v\n", err)
		return
	}

	log.Printf("inserted %d rows\n", res.RowsAffected())

	// Делаем запрос на выборку записей из таблицы note
	rows, err := con.Query(ctx, "SELECT id, title, body, created_at, updated_at FROM note")
	if err != nil {
		log.Printf("failed to select notes: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, body string
		var createdAt time.Time
		var updatedAt sql.NullTime

		err = rows.Scan(&id, &title, &body, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("failed to scan note: %v\n", err)
			return
		}

		log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updatedAt)
	}
}
