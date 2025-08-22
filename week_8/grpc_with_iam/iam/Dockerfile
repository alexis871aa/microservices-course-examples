# ============================
# Stage 1: Build stage
# ============================

# Используем официальный образ Go на базе Alpine — лёгкий, быстрый и безопасный
FROM golang:1.24.2-alpine AS builder

# Устанавливаем git — нужен для загрузки зависимостей из приватных и публичных репозиториев
RUN apk add --no-cache git

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем workspace файлы и модули — это позволяет Go правильно загрузить все зависимости
COPY go.work ./
COPY shared/ ./shared/
COPY ufo/go.mod ufo/go.sum ./ufo/
COPY iam/go.mod iam/go.sum ./iam/

# Переходим в директорию iam для работы
WORKDIR /app/iam

# Загружаем все зависимости, указанные в модульных файлах
# Делаем это перед копированием исходного кода для эффективного кэширования Docker слоёв
RUN go mod download

# Копируем весь исходный код iam сервиса
COPY iam/ ./

# Скачиваем grpc-health-probe — это утилита для проверки состояния gRPC-сервиса
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe

# Делаем скачанный файл исполняемым
RUN chmod +x grpc-health-probe

# Собираем бинарный файл IAM gRPC сервиса для Linux-архитектуры без CGO
# CGO_ENABLED=0 — отключаем CGO для статической линковки (нет зависимостей от libc)
# GOOS=linux GOARCH=amd64 — явно указываем целевую платформу
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o iam-grpc-server ./cmd/grpc_server/main.go

# ============================
# Stage 2: Final image
# ============================

# Используем чистый Alpine образ как итоговый контейнер — он будет лёгким и безопасным
FROM alpine:3.21.3

# Создаём системного пользователя без root-прав — best practice для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Запускаем приложение под non-root пользователем
USER appuser

# Копируем скомпилированный бинарник из стадии builder
COPY --from=builder /app/iam/iam-grpc-server ./iam-grpc-server

# Копируем бинарь grpc-health-probe в /bin
COPY --from=builder /app/iam/grpc-health-probe /bin/grpc-health-probe

# Экспонируем порт gRPC сервера
EXPOSE 50052

# Устанавливаем команду запуска — запускаем наш gRPC сервер
ENTRYPOINT ["./iam-grpc-server"] 