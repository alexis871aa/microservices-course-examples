# UFO Sightings gRPC Example с Interceptors

Этот проект представляет собой пример реализации gRPC сервиса для фиксации и управления наблюдениями НЛО с использованием интерцепторов.

### Особенности

- Простой CRUD API с использованием gRPC
- Использование UUID для идентификаторов
- Поддержка нулабельных полей через google.protobuf.wrappers
- Метки времени для создания, обновления и удаления
- Реализация паттерна частичного обновления для метода Update
- **Использование gRPC интерцепторов для:**
  - Логирования запросов
  - Метрик (замер времени выполнения запросов)

### Структура проекта

```
.
├── cmd
│   ├── grpc_client      # gRPC клиент
│   └── grpc_server      # gRPC сервер с интерцепторами
├── internal
│   └── interceptors     # Реализация интерцепторов
├── pkg
│   └── proto            # Сгенерированный Go код из proto-файлов
├── proto
│   ├── buf.gen.yaml     # Конфигурация для генерации кода
│   ├── buf.yaml         # Конфигурация для линтинга proto-файлов
│   └── ufo/v1           # Proto-файлы версии v1
└── go.mod               # Go модуль
```

### Зависимости

- Go 1.24+
- [buf](https://github.com/bufbuild/buf)
- [Task](https://taskfile.dev/)
- [protoc-gen-go](https://developers.google.com/protocol-buffers/docs/reference/go-generated)
- [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc)

### Начало работы

1. Установите [Go](https://golang.org/doc/install)
2. Установите [Task](https://taskfile.dev/#/installation)
3. Установите нужные инструменты:

```bash
task install-buf
task proto:install-plugins
task install-golangci-lint
```

### Генерация кода из proto-файлов

```bash
task proto:gen
```

### Запуск сервера

```bash
go run cmd/grpc_server/main.go
```

### Запуск клиента

```bash
go run cmd/grpc_client/main.go
```

## Интерцепторы

В этом примере реализованы следующие gRPC интерцепторы:

### Серверные интерцепторы

- **LoggingInterceptor**: Логирует информацию о входящих запросах и о времени их выполнения.

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
- Реализованы unary интерцепторы для обработки унарных запросов

## Линтинг

Для запуска линтеров используйте:

```bash
task lint
``` 
