# Weather Bot 🌤️

Простой Telegram бот для получения информации о погоде в любом городе мира с использованием бесплатного Open-Meteo API.

## Возможности

- 🌍 **Погода в любом городе** - `/weather <город>`
- 🏙️ **Популярные города** - быстрый доступ к крупным городам России
- ⚡ **Inline клавиатуры** - удобное взаимодействие через кнопки
- 🎨 **Красивое оформление** - эмодзи в зависимости от типа погоды
- 🔄 **Обновление данных** - кнопка для получения свежих данных
- 🆓 **Бесплатно** - без API ключей и регистрации

## Демонстрируемые технологии

- **Внешние API** - интеграция с Open-Meteo (Geocoding + Weather)
- **HTTP запросы** - работа с REST API
- **JSON парсинг** - обработка ответов API
- **Graceful shutdown** - корректная остановка сервиса
- **Inline клавиатуры** - современный UX в Telegram
- **Структурированный код** - методы на структуре, разделение ответственности
- **WMO коды погоды** - стандартные коды Всемирной метеорологической организации

## Настройка Telegram Bot API

### Регистрация обработчиков команд

В Telegram боте используется система обработчиков для разных типов сообщений:

```go
// Регистрируем обработчики
b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, wb.startHandler)
b.RegisterHandler(bot.HandlerTypeMessageText, "/weather", bot.MatchTypePrefix, wb.weatherHandler)
b.RegisterHandler(bot.HandlerTypeMessageText, "/popular", bot.MatchTypeExact, wb.popularCitiesHandler)
b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix, wb.callbackHandler)
```

**Типы обработчиков:**
- `HandlerTypeMessageText` - обрабатывает текстовые сообщения
- `HandlerTypeCallbackQueryData` - обрабатывает нажатия на inline кнопки

**Типы сопоставления:**
- `MatchTypeExact` - точное совпадение (например, только "/start")
- `MatchTypePrefix` - совпадение по префиксу (например, "/weather Москва")

### Inline клавиатуры и Callback'и

#### Создание inline клавиатуры

```go
keyboard := &models.InlineKeyboardMarkup{
    InlineKeyboard: [][]models.InlineKeyboardButton{
        {
            {Text: "🌍 Популярные города", CallbackData: "popular_cities"},
        },
        {
            {Text: "🔄 Обновить", CallbackData: "weather_" + city},
            {Text: "🌍 Другие города", CallbackData: "popular_cities"},
        },
    },
}
```

**Структура:**
- `InlineKeyboard` - двумерный массив кнопок (ряды и столбцы)
- `Text` - текст, отображаемый на кнопке
- `CallbackData` - данные, отправляемые при нажатии (до 64 байт)

#### Обработка callback'ов

```go
func (wb *WeatherBot) callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
    callback := update.CallbackQuery
    data := callback.Data // Получаем CallbackData

    // Обязательно отвечаем на callback query
    b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
        CallbackQueryID: callback.ID,
    })

    // Получаем chat ID из сообщения
    if callback.Message.Message == nil {
        return // Сообщение недоступно
    }
    chatID := callback.Message.Message.Chat.ID

    // Обрабатываем данные
    if data == "popular_cities" {
        // Показываем популярные города
    } else if strings.HasPrefix(data, "weather_") {
        city := strings.TrimPrefix(data, "weather_")
        // Показываем погоду для города
    }
}
```

**Важные моменты:**
- Всегда вызывайте `AnswerCallbackQuery` - иначе у пользователя будет "крутиться" индикатор загрузки
- `callback.Message` может быть `MaybeInaccessibleMessage` - проверяйте доступность
- `CallbackData` ограничена 64 байтами - используйте короткие идентификаторы

#### Редактирование сообщений

```go
// Редактируем текст сообщения с новой клавиатурой
b.EditMessageText(ctx, &bot.EditMessageTextParams{
    ChatID:      chatID,
    MessageID:   loadingMsg.ID,
    Text:        weatherText,
    ReplyMarkup: keyboard,
})
```

Это позволяет обновлять содержимое сообщения без создания нового.

## Исследование внешних API

### Как изучить API с нуля

#### 1. Документация API

**Open-Meteo Geocoding API:**
- Документация: https://open-meteo.com/en/docs/geocoding-api
- Endpoint: `https://geocoding-api.open-meteo.com/v1/search`

**Open-Meteo Weather API:**
- Документация: https://open-meteo.com/en/docs
- Endpoint: `https://api.open-meteo.com/v1/forecast`

#### 2. Тестирование через curl

```bash
# Поиск координат города
curl "https://geocoding-api.open-meteo.com/v1/search?name=Moscow&count=1&language=ru&format=json"

# Получение погоды по координатам
curl "https://api.open-meteo.com/v1/forecast?latitude=55.7558&longitude=37.6176&current=temperature_2m,apparent_temperature,relative_humidity_2m,weather_code,surface_pressure,wind_speed_10m,wind_direction_10m&timezone=auto"
```

#### 3. Анализ JSON ответов

**Geocoding API ответ:**
```json
{
  "results": [
    {
      "id": 524901,
      "name": "Moscow",
      "latitude": 55.75222,
      "longitude": 37.61556,
      "elevation": 144.0,
      "feature_code": "PPLC",
      "country_code": "RU",
      "admin1_id": 524894,
      "admin2_id": 524901,
      "admin3_id": 524901,
      "timezone": "Europe/Moscow",
      "population": 10381222,
      "country_id": 2017370,
      "country": "Russia",
      "admin1": "Moscow",
      "admin2": "Moscow",
      "admin3": "Moscow"
    }
  ],
  "generationtime_ms": 1.2345
}
```

**Weather API ответ:**
```json
{
  "latitude": 55.75,
  "longitude": 37.625,
  "generationtime_ms": 0.123,
  "utc_offset_seconds": 10800,
  "timezone": "Europe/Moscow",
  "timezone_abbreviation": "MSK",
  "elevation": 144.0,
  "current_units": {
    "time": "iso8601",
    "interval": "seconds",
    "temperature_2m": "°C",
    "apparent_temperature": "°C",
    "relative_humidity_2m": "%",
    "weather_code": "wmo code",
    "surface_pressure": "hPa",
    "wind_speed_10m": "m/s",
    "wind_direction_10m": "°"
  },
  "current": {
    "time": "2024-01-15T12:00",
    "interval": 900,
    "temperature_2m": -5.2,
    "apparent_temperature": -8.1,
    "relative_humidity_2m": 87,
    "weather_code": 3,
    "surface_pressure": 1013.2,
    "wind_speed_10m": 3.4,
    "wind_direction_10m": 245
  }
}
```

#### 4. Создание Go структур

На основе JSON ответов создаем Go структуры:

```go
// Для Geocoding API
type GeocodingResponse []struct {
    Name      string  `json:"name"`
    Country   string  `json:"country"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

// Для Weather API
type WeatherResponse struct {
    Current struct {
        Time                string  `json:"time"`
        Temperature2m       float64 `json:"temperature_2m"`
        ApparentTemperature float64 `json:"apparent_temperature"`
        RelativeHumidity2m  int     `json:"relative_humidity_2m"`
        WeatherCode         int     `json:"weather_code"`
        SurfacePressure     float64 `json:"surface_pressure"`
        WindSpeed10m        float64 `json:"wind_speed_10m"`
        WindDirection10m    float64 `json:"wind_direction_10m"`
    } `json:"current"`
}
```

#### 5. Инструменты для исследования API

**Онлайн инструменты:**
- **Postman** - GUI для тестирования API
- **Insomnia** - альтернатива Postman
- **HTTPie** - командная строка для HTTP запросов
- **JSON Formatter** - форматирование и анализ JSON

**Браузерные инструменты:**
- Просто вставьте URL в браузер для GET запросов
- Developer Tools (F12) → Network для анализа запросов

**Go инструменты:**
```go
// Для отладки - выводим сырой JSON
var raw json.RawMessage
json.NewDecoder(resp.Body).Decode(&raw)
fmt.Printf("Raw JSON: %s\n", raw)
```

#### 6. Обработка ошибок API

```go
// Проверяем HTTP статус
if resp.StatusCode != http.StatusOK {
    body, _ := io.ReadAll(resp.Body)
    return fmt.Errorf("API error %d: %s", resp.StatusCode, body)
}

// Проверяем структуру ответа
if len(geocoding.Results) == 0 {
    return fmt.Errorf("city not found")
}
```

### WMO Weather Codes

Open-Meteo использует стандартные коды WMO (Всемирная метеорологическая организация):

| Код | Описание |
|-----|----------|
| 0 | Ясно |
| 1 | Преимущественно ясно |
| 2 | Переменная облачность |
| 3 | Пасмурно |
| 45, 48 | Туман |
| 51, 53, 55 | Морось |
| 61, 63, 65 | Дождь |
| 71, 73, 75 | Снег |
| 95 | Гроза |
| 96, 99 | Гроза с градом |

Полная таблица: https://open-meteo.com/en/docs#weathervariables

## Настройка

### 1. Создание Telegram бота
1. Найдите @BotFather в Telegram
2. Отправьте `/newbot`
3. Следуйте инструкциям
4. Сохраните токен

### 2. Настройка токена в коде
Откройте `cmd/main.go` и замените константу:

```go
const (
    // Замените на ваш токен от @BotFather
    telegramBotToken = "1234567890:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
)
```

**Примечание:** API ключ для погоды не нужен! Open-Meteo предоставляет бесплатный доступ без регистрации.

### 3. Установка зависимостей
```bash
cd week_5/telegram/weather-bot
go mod tidy
```

## Запуск

```bash
# Из директории weather-bot
go run cmd/main.go

# Или соберите и запустите
go build -o weather-bot cmd/main.go
./weather-bot
```

## Использование

### Команды
- `/start` - показать приветствие и помощь
- `/weather <город>` - получить погоду для указанного города
- `/popular` - показать кнопки с популярными городами

### Примеры
```
/weather Москва
/weather London
/weather New York
/weather Токио
```

### Inline кнопки
- 🌍 **Популярные города** - быстрый выбор из списка
- 🔄 **Обновить** - получить свежие данные о погоде
- 🏙️ **Города** - кнопки для каждого популярного города

## Архитектура

```
WeatherBot
├── Telegram Bot API (go-telegram/bot)
├── Open-Meteo Geocoding API (поиск координат города)
├── Open-Meteo Weather API (данные о погоде)
└── Graceful Shutdown
```

### Основные компоненты

1. **WeatherBot struct** - основная структура бота
2. **Geocoding Integration** - поиск координат города через Open-Meteo
3. **Weather API Integration** - получение данных о погоде
4. **Message Handlers** - обработчики команд и callback'ов
5. **UI Components** - inline клавиатуры и форматирование сообщений
6. **WMO Weather Codes** - интерпретация стандартных кодов погоды

### API Endpoints

- **Geocoding**: `https://geocoding-api.open-meteo.com/v1/search`
- **Weather**: `https://api.open-meteo.com/v1/forecast`

### Обработка ошибок
- Информативные сообщения об ошибках API
- Graceful degradation при недоступности сервиса
- Проверка корректности токена при запуске
- Обработка случаев, когда город не найден

## Данные о погоде

Бот показывает:
- 🌡️ **Температура** (текущая и ощущаемая)
- 📝 **Описание** (на основе WMO кодов)
- 💨 **Ветер** (скорость и направление)
- 💧 **Влажность**
- 📊 **Давление**
- 🕐 **Время обновления**

## Логи

Бот выводит логи:
- Поиск координат городов
- Получение данных от API
- Ошибки HTTP запросов
- Graceful shutdown

## Остановка

Нажмите `Ctrl+C` для graceful shutdown бота.

## Преимущества Open-Meteo

- ✅ **Бесплатно** - без API ключей и лимитов
- ✅ **Без регистрации** - сразу готов к использованию
- ✅ **Высокая точность** - данные от национальных метеослужб
- ✅ **Открытый исходный код** - прозрачность алгоритмов
- ✅ **Быстрые обновления** - данные обновляются каждый час

## Расширения

Легко добавить:
- Кэширование ответов (Redis/in-memory)
- Прогноз на несколько дней (Open-Meteo поддерживает до 16 дней)
- Исторические данные (до 80 лет назад)
- Графики погоды
- Уведомления о погоде
- Геолокацию пользователя
- Избранные города
- Переменные окружения для токенов
- Дополнительные данные (UV индекс, качество воздуха, пыльца) 