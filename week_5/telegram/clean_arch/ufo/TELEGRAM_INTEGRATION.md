# Telegram Integration для UFO Service

Интеграция с Telegram Bot API для отправки уведомлений о новых наблюдениях UFO.

## Архитектура

```
UFO Service
├── internal/client/http/telegram.go     - HTTP клиент для Telegram Bot API
├── internal/service/telegram/service.go - Сервис для работы с уведомлениями
├── internal/service/ufo/create.go       - Интеграция в UFO сервис
└── internal/app/di.go                   - DI контейнер
```

## Компоненты

### 1. HTTP Клиент (`internal/client/http/telegram.go`)

**Интерфейс:**
```go
type TelegramClient interface {
    SendMessage(ctx context.Context, chatID int64, text string) error
}
```

**Реализация:**
- Использует библиотеку `github.com/go-telegram/bot`
- Отправляет текстовые сообщения в указанный чат
- Логирует успешную отправку

### 2. Telegram Сервис (`internal/service/telegram/service.go`)

**Интерфейс:**
```go
type Service interface {
    SendUFONotification(ctx context.Context, uuid string, sighting model.SightingInfo) error
}
```

**Функциональность:**
- Формирует сообщение из шаблона с данными UFO
- Отправляет уведомление через HTTP клиент
- Обрабатывает все поля модели `SightingInfo`

**Шаблон сообщения:**
```
🛸 **НОВОЕ НАБЛЮДЕНИЕ UFO!**

🆔 **ID:** `uuid`
📍 **Место:** location
📝 **Описание:** description
🕐 **Время наблюдения:** observed_at (если указано)
🎨 **Цвет:** color (если указан)
🔊 **Звук:** Да/Нет (если указано)
⏱️ **Длительность:** duration сек (если указана)

📅 **Зарегистрировано:** current_timestamp
```

### 3. Интеграция в UFO Сервис

В методе `Create` UFO сервиса добавлен вызов Telegram уведомления:

```go
func (s *service) Create(ctx context.Context, info model.SightingInfo) (string, error) {
    uuid, err := s.ufoRepository.Create(ctx, info)
    if err != nil {
        return "", err
    }

    // Отправляем уведомление в Telegram
    if err := s.telegramService.SendUFONotification(ctx, uuid, info); err != nil {
        // Логируем ошибку, но не прерываем выполнение
        log.Printf("Failed to send telegram notification for UFO %s: %v", uuid, err)
    }

    return uuid, nil
}
```

**Особенности:**
- Ошибка отправки уведомления не прерывает создание UFO
- Ошибки логируются для мониторинга
- Асинхронная отправка (можно легко добавить)

### 4. DI Контейнер (`internal/app/di.go`)

**Новые зависимости:**
```go
type diContainer struct {
    // ... существующие поля
    telegramService service.TelegramService
    telegramClient  httpClient.TelegramClient
}
```

**Методы:**
- `TelegramClient(ctx)` - создает HTTP клиент с токеном
- `TelegramService(ctx)` - создает сервис с клиентом и chat ID
- `PartService(ctx)` - обновлен для принятия Telegram сервиса

## Конфигурация

В `internal/app/di.go` захардкожены константы:

```go
const (
    telegramBotToken = "YOUR_TELEGRAM_BOT_TOKEN_HERE"
    telegramChatID   = -1001234567890 // Замените на реальный chat ID
)
```

### Получение токена бота:
1. Найдите @BotFather в Telegram
2. Отправьте `/newbot`
3. Следуйте инструкциям
4. Сохраните токен

### Получение chat ID:

**Способ 1: Через веб-интерфейс Telegram**
1. Откройте Telegram Web (web.telegram.org)
2. Выберите чат/группу/канал
3. Посмотрите в URL: `https://web.telegram.org/k/#-1001234567890`
4. Число после `#` - это chat ID

**Способ 2: Через Bot API**
1. Добавьте бота в группу/канал или напишите ему в личку
2. Отправьте любое сообщение
3. Выполните запрос: `https://api.telegram.org/bot{TOKEN}/getUpdates`
4. Найдите `chat.id` в ответе

**Способ 3: Через @userinfobot**
1. Найдите @userinfobot в Telegram
2. Добавьте его в группу или напишите `/start`
3. Бот покажет chat ID

**Способ 4: Через curl**
```bash
curl "https://api.telegram.org/bot{YOUR_TOKEN}/getUpdates"
```

**Примеры chat ID:**
- Личный чат: `123456789` (положительное число)
- Группа: `-123456789` (отрицательное число)
- Супергруппа/канал: `-1001234567890` (начинается с -100)

**Для тестирования:**
Можете использовать свой личный chat ID - просто напишите боту `/start` и получите updates.

## Зависимости

Добавлена зависимость в `go.mod`:
```
github.com/go-telegram/bot v1.15.0
```

## Использование

1. **Настройте токен и chat ID** в `di.go`
2. **Запустите сервис**
3. **Создайте UFO через API** - уведомление отправится автоматически

Пример создания UFO:
```bash
grpcurl -plaintext -d '{
  "info": {
    "location": "Москва, Красная площадь",
    "description": "Яркий объект в форме диска",
    "color": "Серебристый",
    "sound": true,
    "duration_seconds": 120
  }
}' localhost:50051 ufo.v1.UFOService/Create
```

## Мониторинг

Логи Telegram интеграции:
- Успешная отправка: `Telegram message sent to chat {chatID}`
- Ошибка отправки: `Failed to send telegram notification for UFO {uuid}: {error}`

## Расширения

Легко добавить:
- **Переменные окружения** для токена и chat ID
- **Асинхронную отправку** через очереди
- **Retry механизм** при ошибках
- **Форматирование Markdown** для красивых сообщений
- **Inline кнопки** для взаимодействия с UFO
- **Уведомления об обновлениях** UFO
- **Множественные чаты** для разных типов уведомлений 