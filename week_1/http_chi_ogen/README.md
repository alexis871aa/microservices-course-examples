# HTTP-сервис с Chi роутером и Ogen

Простой HTTP-сервис для работы с данными о погоде, построенный на базе роутера [Chi](https://github.com/go-chi/chi) и генератора кода [Ogen](https://github.com/ogen-go/ogen) на основе OpenAPI спецификации.

## Структура проекта

```
.
├── api/
│   └── openapi/        # OpenAPI спецификации
│       └── weather.yaml # Спецификация API погоды
├── cmd/
│   ├── http_client/    # HTTP-клиент для тестирования API
│   └── http_server/    # HTTP-сервер на базе Chi роутера
├── internal/
│   └── middleware/      # Middleware для сервера
├── pkg/
│   └── weatherapi/     # Сгенерированный API код из OpenAPI спецификации
├── .golangci.yml       # Конфигурация линтера
├── go.mod              # Зависимости Go-модуля
└── Taskfile.yaml       # Задачи для управления проектом
```

## Особенности

- CRUD API для работы с данными о погоде
- Автоматическая генерация кода сервера из OpenAPI спецификации
- Типизированные обработчики благодаря Ogen
- Удобная маршрутизация с помощью Chi
- Автоматическая валидация запросов на основе схем OpenAPI

## Запуск проекта

### Генерация кода из OpenAPI спецификации

```bash
task ogen:gen
```

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
```json
{
  "message": "Weather for city 'SomeCity' not found",
  "code": "not_found"
}
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