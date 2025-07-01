# UFO Sightings gRPC Example

Этот проект представляет собой пример реализации gRPC сервиса для фиксации и управления наблюдениями НЛО.

## О gRPC

[gRPC](https://grpc.io/) - это высокопроизводительный фреймворк RPC (Remote Procedure Call), разработанный Google, который позволяет клиентам и серверам общаться прозрачно и упрощает создание распределенных систем.

Основные преимущества:
- Высокая производительность
- Строгая типизация (с использованием Protocol Buffers)
- Независимость от языка программирования
- Поддержка потоковой передачи данных
- Встроенный механизм генерации клиентского и серверного кода

### Особенности

- Простой CRUD API с использованием gRPC
- Использование UUID для идентификаторов
- Поддержка нулабельных полей через google.protobuf.wrappers
- Метки времени для создания, обновления и удаления
- Реализация паттерна частичного обновления для метода Update

### Структура проекта

```
.
├── cmd/
│   ├── grpc_client/     # gRPC клиент
│   └── grpc_server/     # gRPC сервер
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

### Зависимости

- Go 1.24+
- [buf](https://github.com/bufbuild/buf) - современная система для работы с Protocol Buffers
- [Task](https://taskfile.dev/) - альтернатива Make для запуска задач
- [protoc-gen-go](https://developers.google.com/protocol-buffers/docs/reference/go-generated) - генератор Go кода из proto файлов
- [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc) - генератор gRPC кода для Go

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

Сервер запустится на порту 50051.

### Запуск клиента

```bash
go run cmd/grpc_client/main.go
```

Клиент подключится к серверу и выполнит несколько тестовых операций.

## API методы

gRPC сервис `UFOService` предоставляет следующие методы:

### Create
Создает новое наблюдение НЛО.

### Get
Получает информацию о наблюдении НЛО по UUID.

### Update
Обновляет существующее наблюдение НЛО.

### Delete
Удаляет наблюдение НЛО (мягкое удаление).

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
- Graceful shutdown для корректного завершения работы сервера

## Линтинг

Для запуска линтеров используйте:

```bash
task lint
``` 
