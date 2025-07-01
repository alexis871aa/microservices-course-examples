# HTTP-сервис с Chi роутером

Простой HTTP-сервис для работы с данными о погоде, построенный на базе роутера [Chi](https://github.com/go-chi/chi).

## Структура проекта

```
.
├── cmd/
│   ├── http_client/    # HTTP-клиент для тестирования API
│   └── http_server/    # HTTP-сервер на базе Chi роутера
├── pkg/
│   └── models/         # Модели данных
├── .golangci.yml       # Конфигурация линтера
├── go.mod              # Зависимости Go-модуля
└── Taskfile.yaml       # Задачи для управления проектом
```

## Особенности

- REST API для работы с данными о погоде
- Использование [Chi](https://github.com/go-chi/chi) для маршрутизации HTTP запросов
- Реализация middleware для:
  - Логирования запросов/ответов
  - Обработки ошибок
  - Установки заголовков Content-Type
  - Сжатия ответов (при необходимости)
- Структурированное логирование
- Graceful shutdown сервера
- Модульная архитектура с разделением на слои

## Запуск проекта

### Запуск сервера

```bash
go run cmd/http_server/main.go
```

### Запуск клиента

```bash
go run cmd/http_client/main.go
```

## Описание API

### GET /api/weather/{city}

Получение данных о погоде для указанного города.

**Ответ (200 OK)**:
```json
{
  "city": "Moscow",
  "temperature": 25.5,
  "updated_at": "2023-05-15T10:30:00Z"
}
```

**Ответ (404 Not Found)**:
```
Weather for city 'SomeCity' not found
```

### PUT /api/weather/{city}

Обновление данных о погоде для указанного города.

**Запрос**:
```json
{
  "temperature": 25.5
}
```

**Ответ (200 OK)**:
```json
{
  "city": "Moscow",
  "temperature": 25.5,
  "updated_at": "2023-05-15T10:30:00Z"
}
```

### POST /api/weather

Создание данных о погоде для нового города.

**Запрос**:
```json
{
  "city": "Berlin",
  "temperature": 20.5
}
```

**Ответ (201 Created)**:
```json
{
  "city": "Berlin",
  "temperature": 20.5,
  "updated_at": "2023-05-15T10:30:00Z"
}
```

### DELETE /api/weather/{city}

Удаление данных о погоде для указанного города.

**Ответ (204 No Content)**

## Линтинг

Для запуска линтеров используйте:

```bash
task lint
``` 