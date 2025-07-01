# JWT Service

Демонстрационный gRPC сервис для работы с JWT токенами.

## Возможности

- **Login** - аутентификация пользователя и получение пары токенов
- **GetAccessToken** - получение нового access токена по refresh токену  
- **GetRefreshToken** - получение нового refresh токена

## Хардкодные пользователи

Пароли хранятся в виде bcrypt хешей, но для логина используйте обычные пароли:

- `admin:admin123`
- `user1:password1`
- `user2:password2`
- `john:john123`
- `alice:alice456`

## Запуск

```bash
# Генерация proto файлов
task proto:gen

# Запуск сервера
go run cmd/server/main.go

# Или через task
task run
```

Сервер запустится на порту `:50051`.

## Тестирование

### Автоматическое тестирование
```bash
# Полный набор тестов (позитивные + негативные)
task test:all

# Только позитивные сценарии
task test

# Тест логина для всех пользователей
task test:login

# Тест негативных сценариев (ошибки)
task test:fail
```

### Ручное тестирование с grpcurl

#### Логин
```bash
grpcurl -plaintext -d '{
  "username": "admin",
  "password": "admin123"
}' localhost:50051 jwt.v1.JWTService/Login
```

#### Получение нового access токена
```bash
grpcurl -plaintext -d '{
  "refresh_token": "YOUR_REFRESH_TOKEN"
}' localhost:50051 jwt.v1.JWTService/GetAccessToken
```

#### Получение нового refresh токена
```bash
grpcurl -plaintext -d '{
  "refresh_token": "YOUR_REFRESH_TOKEN"
}' localhost:50051 jwt.v1.JWTService/GetRefreshToken
```

## Структура проекта

```
├── proto/jwt/v1/           # Proto файлы
├── pkg/proto/jwt/v1/       # Сгенерированный Go код
├── internal/
│   ├── model/              # Модели данных
│   ├── service/            # Бизнес-логика JWT
│   │   ├── jwt.go          # Основная структура и конструктор
│   │   ├── auth.go         # Логика аутентификации
│   │   ├── token_generator.go  # Генерация токенов
│   │   └── token_validator.go  # Валидация токенов
│   └── api/                # gRPC хендлеры
├── cmd/server/             # Точка входа
└── README.md
```

## Тестовые сценарии

### ✅ Позитивные сценарии
- Логин всех 5 пользователей
- Генерация пары токенов
- Обновление access токена по refresh
- Обновление refresh токена

### ❌ Негативные сценарии  
- Неправильный пароль
- Несуществующий пользователь
- Пустые поля логина
- Невалидный refresh токен
- Пустой refresh токен

## Особенности

- Access токены живут 15 минут
- Refresh токены живут 24 часа
- Используется HMAC-SHA256 для подписи
- Пароли хранятся как bcrypt хеши
- Простая архитектура без слоев (демонстрационный пример)
- Логика разбита по файлам для удобства чтения
- Полное покрытие тестами позитивных и негативных сценариев 