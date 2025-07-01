# 🚀 Полный практический гайд по Docker Compose: автоматизируем запуск приложения и базы данных

Этот гайд поможет вам собрать полноценный проект с Docker Compose:  
поднять базу данных Postgres, приложение на Go, настроить миграции через Goose, и всё это запускать одной командой.

Вы научитесь:
- Понимать, зачем нужен Docker Compose и почему одного Docker мало.
- Создавать конфигурацию с переменными окружения.
- Писать миграции с помощью Goose.
- Настраивать автоматическое применение миграций при старте приложения.

📦 Всё оформлено как рабочий проект!

---

## 📖 Зачем нужен Docker Compose?

Допустим, у вас несколько контейнеров:
- Backend-приложение на Go
- База данных PostgreSQL
- (В будущем: PgAdmin, Redis, Kafka и другие)

Без Docker Compose вам придётся вручную запускать каждый контейнер, прописывать порты, сети, volume и переменные окружения.  
**Это неудобно и не масштабируется.**

Docker Compose решает проблему:  
✅ Описываем все сервисы в одном `docker-compose.yml`  
✅ Управляем всем проектом одной командой:
```bash
docker-compose up
```
✅ Добавляем переменные через `.env`  
✅ Удобно масштабируем и поддерживаем проект.

---

## 📂 Структура проекта

```
docker-compose-tutorial
├── cmd
│   └── main.go
├── internal
│   └── migrator
│       └── migrator.go
├── migrations
│   └── 001_create_user_table.sql
├── .env
├── .dockerignore
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

---

## 🚀 Шаг за шагом

### 1. Правим приложение, чтобы порт и путь к миграциям брались из `.env`

В Go-приложении мы подключаемся к базе и указываем директорию миграций через переменные окружения:

- Строка подключения: `DB_URI`
- Путь к миграциям: `MIGRATIONS_DIR`

`.env` файл:
```env
POSTGRES_USER=demo
POSTGRES_PASSWORD=demo
POSTGRES_PORT=5432
APP_PORT=8080
DB_URI=postgres://demo:demo@mypostgres:${POSTGRES_PORT}/postgres
MIGRATIONS_DIR=./migrations
```

---

### 2. Создаём Docker Compose файл

Файл `docker-compose.yml`:

```yaml
version: "3.9"

services:
  postgres:
    image: postgres:15
    container_name: mypostgres
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports:
      - "${POSTGRES_PORT}:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: [ "CMD-SHELL", "pg_isready -U ${POSTGRES_USER}" ]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    networks:
      - app-network

  app:
    build:
      context: docker-compose-tutorial
      dockerfile: Dockerfile
    container_name: my-go-server
    ports:
      - "${APP_PORT}:8080"
    environment:
      - DB_URI=${DB_URI}
      - MIGRATIONS_DIR=${MIGRATIONS_DIR}
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped
    networks:
      - app-network

networks:
  app-network:
    driver: bridge

volumes:
  pgdata:
```

---

### 3. Создаём файл `.env`

Как указано выше, файл `.env` содержит все параметры окружения.

---

### 4. Запускаем контейнеры

```bash
docker-compose up --build
```

Проверяем логи:
```bash
docker-compose logs -f
```

Приложение поднимется на [http://localhost:8080](http://localhost:8080)

---

### 5. Понимаем, что нет таблицы

Попробуем отправить запрос:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "email": "alice@example.com"}'
```

И видим ошибку: **relation "users" does not exist**  
Нам нужно применить миграции!

---

### 6. Добавляем миграции с помощью Goose

🦢 **Goose** — это утилита для управления миграциями базы данных.

Файл миграции `migrations/001_create_user_table.sql`:

```sql
-- +goose Up
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username text,
  email text,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS users;
```

---

### 7. Пишем мигратор и пробрасываем миграции в Dockerfile

Финальный `migrator.go`:

```go
package migrator

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Up() error {
	return goose.Up(m.db, m.migrationsDir)
}
```

В Dockerfile добавляем:

```Dockerfile
# Копируем папку с миграциями
COPY --from=builder /app/migrations ./migrations
```

И не забываем прокинуть путь в переменную `MIGRATIONS_DIR` через `.env` и `docker-compose.yml`.

---

### 8. Ребилдим и проверяем

Пересобираем контейнеры с нуля:

```bash
docker-compose build --no-cache
```

Запускаем:

```bash
docker-compose up -d
```

Проверяем:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "email": "alice@example.com"}'
```

🎉 Успех: **"Пользователь успешно создан"**

---

## 🎯 Итог

Теперь у нас:
- Автоматизированный запуск приложения и базы данных.
- Автоматический прогон миграций при старте.
- Простая структура проекта для изучения и масштабирования.

---

## 📬 Контакты автора

Если вам понравился урок — заглядывайте:
- [Telegram канал](https://t.me/olezhek28go) — полезные заметки, лайф стайл и обсуждение волнующих вопросов в мире IT.
- [YouTube канал](https://www.youtube.com/@olezhek28go) — технические и софтовые видео про разработку.

Автор: **Олег Козырев**

---
