# Workspace для микросервисов - Week 3

Данный workspace демонстрирует продвинутую структуру проекта с микросервисами, включающую автоматизацию сборки, деплой через Docker Compose и полную интеграцию инструментов разработки.

## 🏗️ Архитектура проекта

Workspace состоит из следующих компонентов:

### Микросервисы

1. **grpc** - gRPC сервер для работы с наблюдениями НЛО:
   - Полная реализация CRUD операций
   - Protocol Buffers для определения API
   - Работа с UUID, Timestamp и Nullable полями
   - Поддержка gRPC рефлексии для отладки

2. **http** - HTTP сервер на базе Chi и Ogen для работы с данными о погоде:
   - REST API с автогенерацией из OpenAPI спецификации
   - Типизированные обработчики благодаря Ogen
   - Валидация запросов и ответов
   - Swagger UI для документации API

3. **shared** - общие компоненты:
   - `proto/` - Protocol Buffers спецификации для gRPC API
   - `api/` - OpenAPI спецификации для HTTP API
   - `pkg/` - сгенерированный код из спецификаций

### Инфраструктура

4. **deploy/compose/** - Docker Compose конфигурации:
   - `core/` - базовая сетевая инфраструктура
   - `grpc/` - деплой зависимостей gRPC сервиса
   - `http/` - деплой зависимостей HTTP сервиса

## 🚀 Быстрый старт

### Установка зависимостей

```bash
# Установка Task (если не установлен)
go install github.com/go-task/task/v3/cmd/task@latest

# Установка всех инструментов разработки
task install-formatters
task install-golangci-lint
task install-buf
task proto:install-plugins
task ogen:install
task redocly-cli:install
```

### Генерация кода

```bash
# Генерация всего кода из спецификаций
task gen

# Или по отдельности:
task proto:gen    # Генерация из Protocol Buffers
task ogen:gen     # Генерация из OpenAPI
```

### Разработка

```bash
# Форматирование кода
task format

# Линтинг всех модулей
task lint

# Обновление зависимостей
task deps:update
```

## 🐳 Деплой через Docker Compose

### Запуск зависимостей всех сервисов

```bash
# Поднять зависимости всех сервисов
task up-all

# Остановить зависимости всех сервисов
task down-all
```

### Запуск отдельных сервисов

```bash
# Базовая инфраструктура (сеть)
task up-core
task down-core

# Зависимости gRPC сервиса
task up-grpc
task down-grpc

# Зависимости HTTP сервиса  
task up-http
task down-http
```

## 🛠️ Локальная разработка

### gRPC сервер

```bash
cd grpc
go run cmd/server/main.go
```

gRPC сервер будет доступен на порту **50051**.

Для тестирования можно использовать:
- [grpcurl](https://github.com/fullstorydev/grpcurl) для CLI
- [BloomRPC](https://github.com/bloomrpc/bloomrpc) для GUI
- [Postman](https://www.postman.com/) с поддержкой gRPC

### HTTP сервер

```bash
cd http
go run cmd/server/main.go
```

HTTP сервер будет доступен на порту **8080**.

Swagger UI доступен по адресу: `http://localhost:8080/swagger`

## 📁 Структура проекта

```
week_3/workspace/
├── grpc/                    # gRPC микросервис
│   ├── cmd/server/         # Точка входа сервера
│   ├── go.mod              # Зависимости модуля
│   └── go.sum
├── http/                    # HTTP микросервис  
│   ├── cmd/server/         # Точка входа сервера
│   ├── go.mod              # Зависимости модуля
│   └── go.sum
├── shared/                  # Общие компоненты
│   ├── api/                # OpenAPI спецификации
│   ├── proto/              # Protocol Buffers
│   ├── pkg/                # Сгенерированный код
│   ├── go.mod              # Зависимости модуля
│   └── go.sum
├── deploy/compose/          # Docker Compose конфигурации
│   ├── core/               # Базовая инфраструктура
│   ├── grpc/               # Деплой gRPC сервиса
│   └── http/               # Деплой HTTP сервиса
├── bin/                     # Локальные инструменты
├── node_modules/            # Node.js зависимости
├── Taskfile.yaml           # Автоматизация задач
├── go.work                 # Go workspace конфигурация
├── buf.work.yaml           # Buf workspace конфигурация
├── .golangci.yml           # Конфигурация линтера
├── package.json            # Node.js зависимости
└── README.md               # Этот файл
```

## 🔧 Инструменты разработки

### Автоматизация (Task)

Все задачи автоматизированы через [Task](https://taskfile.dev/):

- `task gen` - генерация кода из спецификаций
- `task format` - форматирование кода (gofumpt + gci)
- `task lint` - линтинг всех модулей (golangci-lint)
- `task deps:update` - обновление зависимостей
- `task up-all` / `task down-all` - управление Docker Compose

### Генерация кода

- **Protocol Buffers**: [Buf](https://buf.build/) для линтинга и генерации
- **OpenAPI**: [Ogen](https://github.com/ogen-go/ogen) для типизированной генерации
- **Bundling**: [Redocly CLI](https://redocly.com/docs/cli/) для сборки OpenAPI схем

### Качество кода

- **Линтинг**: [golangci-lint](https://golangci-lint.run/) с настроенными правилами
- **Форматирование**: [gofumpt](https://github.com/mvdan/gofumpt) + [gci](https://github.com/daixiang0/gci)
- **Go workspace**: управление зависимостями между модулями

## 🌐 API документация

### gRPC API

Protocol Buffers спецификации находятся в `shared/proto/`. 
Для просмотра доступных методов используйте gRPC рефлексию:

```bash
grpcurl -plaintext localhost:50051 list
```

### HTTP API

OpenAPI спецификации находятся в `shared/api/`.
Swagger UI доступен по адресу: `http://localhost:8080/swagger`

## 🔗 Полезные ссылки

- [gRPC документация](https://grpc.io/docs/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [OpenAPI спецификация](https://spec.openapis.org/oas/latest.html)
- [Chi роутер](https://github.com/go-chi/chi)
- [Ogen генератор](https://github.com/ogen-go/ogen)
- [Task автоматизация](https://taskfile.dev/)
- [Buf инструменты](https://buf.build/)
- [Docker Compose](https://docs.docker.com/compose/) 