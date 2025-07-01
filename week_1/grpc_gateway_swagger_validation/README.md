# UFO Sightings gRPC Example с gRPC-Gateway, Swagger и валидацией

Этот проект представляет собой комплексный пример микросервиса, использующего gRPC с HTTP REST API через gRPC-Gateway, интерактивную Swagger документацию и валидацию входных данных.

## О технологиях

- **gRPC**: высокопроизводительный фреймворк RPC от Google для создания распределенных систем.
- **gRPC-Gateway**: прокси-сервер, который преобразует HTTP/JSON запросы в gRPC и обратно.
- **Protocol Buffers**: механизм сериализации структурированных данных от Google.
- **protoc-gen-validate**: расширение Protocol Buffers для валидации полей.
- **Swagger/OpenAPI**: формат документирования API с интерактивным UI для тестирования.

## Особенности

- Простой CRUD API с использованием gRPC
- REST API через gRPC Gateway
- Swagger UI для интерактивной работы с API
- Использование UUID для идентификаторов
- Поддержка нулабельных полей через google.protobuf.wrappers
- Метки времени для создания, обновления и удаления
- Реализация паттерна частичного обновления для метода Update
- Валидация входящих данных с protoc-gen-validate

## Структура проекта

```
.
├── api/
│   ├── swagger.json     # Сгенерированная OpenAPI спецификация
│   └── swagger-ui.html  # HTML страница для Swagger UI
├── cmd/
│   ├── grpc_client/     # gRPC клиент
│   └── grpc_server/     # gRPC сервер с HTTP Gateway и Swagger UI
├── pkg/
│   └── proto/           # Сгенерированный Go код из proto-файлов
├── proto/
│   ├── buf.gen.yaml     # Конфигурация для генерации кода
│   ├── buf.yaml         # Конфигурация для линтинга proto-файлов
│   └── ufo/v1/          # Proto-файлы версии v1
├── .golangci.yml        # Конфигурация линтера
├── go.mod               # Go модуль
└── Taskfile.yaml        # Задачи для управления проектом
```

## Зависимости

- Go 1.24+
- [buf](https://github.com/bufbuild/buf)
- [Task](https://taskfile.dev/)
- [protoc-gen-go](https://developers.google.com/protocol-buffers/docs/reference/go-generated)
- [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc)
- [protoc-gen-grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [protoc-gen-validate](https://github.com/envoyproxy/protoc-gen-validate)
- [protoc-gen-openapiv2](https://github.com/grpc-ecosystem/grpc-gateway/tree/main/protoc-gen-openapiv2)

## Начало работы

1. Установите [Go](https://golang.org/doc/install)
2. Установите [Task](https://taskfile.dev/#/installation)
3. Установите нужные инструменты:

```bash
task install-buf
task proto:install-plugins
task install-golangci-lint
```

## Генерация кода из proto-файлов

```bash
task proto:gen
```

## Запуск сервера

```bash
go run cmd/grpc_server/main.go
```

Сервер запустит:
- gRPC сервер на порту 50051
- HTTP сервер с REST API через gRPC Gateway и Swagger UI на порту 8081 (доступен по адресу http://localhost:8081)

## Запуск клиента

```bash
go run cmd/grpc_client/main.go
```

## Примеры работы с REST API через gRPC Gateway

Благодаря gRPC Gateway можно взаимодействовать с сервисом через HTTP API.

### Создание нового UFO наблюдения

```bash
curl -X POST http://localhost:8081/api/v1/ufo \
  -H "Content-Type: application/json" \
  -d '{
    "info": {
      "observed_at": "2023-08-15T20:30:00Z",
      "location": "Москва, Останкино",
      "description": "Яркий светящийся объект, перемещался зигзагами",
      "color": "зеленый",
      "sound": true,
      "duration_seconds": 120
    }
  }'
```

Ответ будет содержать UUID созданного наблюдения:
```json
{
  "uuid": "67e55044-10b1-4922-9e8a-4f0d3c2b822b"
}
```

### Пример запроса, нарушающего правила валидации

```bash
curl -X POST http://localhost:8081/api/v1/ufo \
  -H "Content-Type: application/json" \
  -d '{
    "info": {
      "observed_at": "2023-08-15T20:30:00Z",
      "location": "Это очень длинное описание локации, которое превышает максимально допустимую длину в 50 символов и поэтому вызовет ошибку валидации",
      "description": "Яркий светящийся объект",
      "color": "зеленый",
      "sound": true,
      "duration_seconds": 120
    }
  }'
```

Ответ будет содержать ошибку валидации:
```json
{
  "code": 3,
  "message": "validation error: invalid CreateRequest.Info: embedded message failed validation | caused by: invalid SightingInfo.Location: value length must be at most 50 runes"
}
```

### Получение наблюдения по UUID

```bash
curl -X GET http://localhost:8081/api/v1/ufo/67e55044-10b1-4922-9e8a-4f0d3c2b822b
```

Ответ будет выглядеть примерно так:
```json
{
  "sighting": {
    "uuid": "67e55044-10b1-4922-9e8a-4f0d3c2b822b",
    "info": {
      "observed_at": "2023-08-15T20:30:00Z",
      "location": "Москва, Останкино",
      "description": "Яркий светящийся объект, перемещался зигзагами",
      "color": "зеленый",
      "sound": true,
      "duration_seconds": 120
    },
    "created_at": "2023-08-15T21:00:00Z",
    "updated_at": "2023-08-15T21:00:00Z"
  }
}
```

### Обновление наблюдения

```bash
curl -X PATCH http://localhost:8081/api/v1/ufo/67e55044-10b1-4922-9e8a-4f0d3c2b822b \
  -H "Content-Type: application/json" \
  -d '{
    "update_info": {
      "description": "Обновленное описание: светящийся диск",
      "color": "красный"
    }
  }'
```

### Удаление наблюдения

```bash
curl -X DELETE http://localhost:8081/api/v1/ufo/67e55044-10b1-4922-9e8a-4f0d3c2b822b
```

## Преимущества использованных подходов

1. **Единая спецификация API**: один proto-файл для gRPC и REST API
2. **Валидация на уровне спецификации**: правила валидации определяются прямо в proto-файле
3. **Интерактивная документация**: Swagger UI позволяет легко тестировать API
4. **Двойной интерфейс**: API доступен как через gRPC, так и через REST
5. **Типизация**: строгая типизация на всех уровнях взаимодействия

## Сущность Sighting (Наблюдение НЛО)

В текущей версии API сущность Sighting имеет следующую структуру:

### Sighting
- uuid (string): уникальный идентификатор наблюдения
- info (SightingInfo): информация о наблюдении (содержит основные поля)
- created_at (google.protobuf.Timestamp): время создания записи
- updated_at (google.protobuf.Timestamp): время обновления записи
- deleted_at (google.protobuf.Timestamp): время удаления записи (опционально)

### SightingInfo
- observed_at (google.protobuf.Timestamp): время наблюдения НЛО
- location (string): место наблюдения
- description (string): описание наблюдаемого объекта
- color (google.protobuf.StringValue): цвет объекта (опционально)
- sound (google.protobuf.BoolValue): признак наличия звука (опционально)
- duration_seconds (google.protobuf.Int32Value): продолжительность наблюдения в секундах (опционально)

### SightingUpdateInfo
Для частичного обновления используется структура SightingUpdateInfo, в которой все поля опциональны:
- observed_at (google.protobuf.Timestamp): время наблюдения (опционально)
- location (google.protobuf.StringValue): место наблюдения (опционально)
- description (google.protobuf.StringValue): описание наблюдаемого объекта (опционально)
- color (google.protobuf.StringValue): цвет объекта (опционально)
- sound (google.protobuf.BoolValue): признак наличия звука (опционально)
- duration_seconds (google.protobuf.Int32Value): продолжительность наблюдения в секундах (опционально)

## Особенности реализации

- Рефлексия включена на сервере для отладки
- Сервер использует in-memory хранилище (обычную карту с сущностями + RWMutex)
- Клиент показывает простые примеры работы с API: создание, получение, обновление, удаление 
- Валидация запросов перед обработкой через protoc-gen-validate
- Graceful shutdown для корректного завершения работы всех серверов

## Линтинг

Для запуска линтеров используйте:

```bash
task lint
``` 
