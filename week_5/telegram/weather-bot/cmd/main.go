package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Конфигурация - замените на свой токен
const (
	// Получить у @BotFather в Telegram
	telegramBotToken = "TELEGRAM_BOT_TOKEN"
)

// GeocodingResponse структура ответа от Open-Meteo Geocoding API
type GeocodingResponse []struct {
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// WeatherResponse структура ответа от Open-Meteo Weather API
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
	CurrentUnits struct {
		Temperature2m       string `json:"temperature_2m"`
		ApparentTemperature string `json:"apparent_temperature"`
		RelativeHumidity2m  string `json:"relative_humidity_2m"`
		SurfacePressure     string `json:"surface_pressure"`
		WindSpeed10m        string `json:"wind_speed_10m"`
		WindDirection10m    string `json:"wind_direction_10m"`
	} `json:"current_units"`
}

// WeatherBot структура бота
type WeatherBot struct {
	bot *bot.Bot
}

// Популярные города для быстрого доступа
var popularCities = []string{
	"Москва", "Санкт-Петербург", "Новосибирск", "Екатеринбург",
	"Казань", "Нижний Новгород", "Челябинск", "Самара",
}

func main() {
	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем бота
	b, err := bot.New(telegramBotToken,
		bot.WithDefaultHandler(defaultHandler),
	)
	if err != nil {
		log.Printf("Failed to create bot: %v\n", err)
		return
	}

	// Создаем экземпляр WeatherBot
	wb := &WeatherBot{
		bot: b,
	}

	// Регистрируем обработчики
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, wb.startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/weather", bot.MatchTypePrefix, wb.weatherHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/popular", bot.MatchTypeExact, wb.popularCitiesHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix, wb.callbackHandler)

	log.Println("🌤️ Weather Bot started successfully!")

	// Обработка сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping bot...")
		cancel()
	}()

	// Запускаем бота
	b.Start(ctx)
	log.Println("Weather Bot stopped")
}

// startHandler обрабатывает команду /start
func (wb *WeatherBot) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🌍 Популярные города", CallbackData: "popular_cities"},
			},
		},
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: `🌤️ Привет! Я Weather Bot!

Я могу показать погоду в любом городе мира.
Данные предоставляются бесплатным API Open-Meteo.

Команды:
• /weather <город> - погода в указанном городе
• /popular - популярные города России

Примеры:
• /weather Москва
• /weather London
• /weather New York`,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		log.Printf("Failed to send start message: %v\n", err)
		return
	}
}

// weatherHandler обрабатывает команду /weather
func (wb *WeatherBot) weatherHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Извлекаем название города из команды
	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) < 2 {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Укажите название города!\nПример: /weather Москва",
		})
		if err != nil {
			log.Printf("Failed to send error message: %v\n", err)
			return
		}
		return
	}

	city := strings.TrimSpace(parts[1])
	wb.sendWeatherInfo(ctx, b, update.Message.Chat.ID, city)
}

// popularCitiesHandler показывает кнопки с популярными городами
func (wb *WeatherBot) popularCitiesHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	wb.sendPopularCitiesKeyboard(ctx, b, update.Message.Chat.ID)
}

// callbackHandler обрабатывает нажатия на inline кнопки
func (wb *WeatherBot) callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	data := callback.Data

	// Отвечаем на callback query
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
	})
	if err != nil {
		log.Printf("Failed to answer callback query: %v\n", err)
		return
	}

	// Получаем chat ID из сообщения
	if callback.Message.Message == nil {
		log.Printf("Callback message is not accessible")
		return
	}
	chatID := callback.Message.Message.Chat.ID

	if data == "popular_cities" {
		wb.sendPopularCitiesKeyboard(ctx, b, chatID)
	} else if strings.HasPrefix(data, "weather_") {
		city := strings.TrimPrefix(data, "weather_")
		wb.sendWeatherInfo(ctx, b, chatID, city)
	}
}

// sendPopularCitiesKeyboard отправляет клавиатуру с популярными городами
func (wb *WeatherBot) sendPopularCitiesKeyboard(ctx context.Context, b *bot.Bot, chatID int64) {
	var keyboard [][]models.InlineKeyboardButton

	// Создаем кнопки по 2 в ряд
	for i := 0; i < len(popularCities); i += 2 {
		var row []models.InlineKeyboardButton

		row = append(row, models.InlineKeyboardButton{
			Text:         "🏙️ " + popularCities[i],
			CallbackData: "weather_" + popularCities[i],
		})

		if i+1 < len(popularCities) {
			row = append(row, models.InlineKeyboardButton{
				Text:         "🏙️ " + popularCities[i+1],
				CallbackData: "weather_" + popularCities[i+1],
			})
		}

		keyboard = append(keyboard, row)
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "🌍 Выберите город:",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	})
	if err != nil {
		log.Printf("Failed to send popular cities keyboard: %v\n", err)
		return
	}
}

// sendWeatherInfo получает и отправляет информацию о погоде
func (wb *WeatherBot) sendWeatherInfo(ctx context.Context, b *bot.Bot, chatID int64, city string) {
	// Сначала отправляем сообщение о загрузке
	loadingMsg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "🔄 Получаю данные о погоде...",
	})
	if err != nil {
		log.Printf("Failed to send loading message: %v", err)
		return
	}

	// Получаем координаты города
	lat, lon, cityName, country, err := wb.getCityCoordinates(city)
	if err != nil {
		// Редактируем сообщение с ошибкой
		_, editErr := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: loadingMsg.ID,
			Text:      fmt.Sprintf("❌ Не удалось найти город '%s'\n\nПроверьте правильность написания города.", city),
		})
		if editErr != nil {
			log.Printf("Failed to edit error message: %v\n", editErr)
			return
		}
		return
	}

	// Получаем данные о погоде
	weather, err := wb.getWeatherData(lat, lon)
	if err != nil {
		// Редактируем сообщение с ошибкой
		_, editErr := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: loadingMsg.ID,
			Text:      fmt.Sprintf("❌ Не удалось получить погоду для города '%s'\n\nПопробуйте позже.", cityName),
		})
		if editErr != nil {
			log.Printf("Failed to edit error message: %v\n", editErr)
			return
		}
		return
	}

	// Формируем красивое сообщение с погодой
	weatherText := wb.formatWeatherMessage(weather, cityName, country)

	// Создаем кнопки для обновления и выбора других городов
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🔄 Обновить", CallbackData: "weather_" + city},
				{Text: "🌍 Другие города", CallbackData: "popular_cities"},
			},
		},
	}

	// Редактируем сообщение с результатом
	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   loadingMsg.ID,
		Text:        weatherText,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		log.Printf("Failed to edit weather message: %v\n", err)
		return
	}
}

// getCityCoordinates получает координаты города через Open-Meteo Geocoding API
func (wb *WeatherBot) getCityCoordinates(city string) (float64, float64, string, string, error) {
	// Формируем URL для запроса к Geocoding API
	baseURL := "https://geocoding-api.open-meteo.com/v1/search"
	params := url.Values{}
	params.Add("name", city)
	params.Add("count", "1")
	params.Add("language", "ru")
	params.Add("format", "json")

	fullURL := baseURL + "?" + params.Encode()

	// Делаем HTTP запрос
	resp, err := http.Get(fullURL)
	if err != nil {
		return 0, 0, "", "", fmt.Errorf("failed to make geocoding request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Failed to close response body: %v\n", cerr)
		}
	}()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return 0, 0, "", "", fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	// Декодируем JSON ответ
	var geocoding struct {
		Results GeocodingResponse `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&geocoding); err != nil {
		return 0, 0, "", "", fmt.Errorf("failed to decode geocoding response: %w", err)
	}

	// Проверяем, что найден хотя бы один результат
	if len(geocoding.Results) == 0 {
		return 0, 0, "", "", fmt.Errorf("city not found")
	}

	result := geocoding.Results[0]
	log.Printf("Found coordinates for %s, %s: %.4f, %.4f\n", result.Name, result.Country, result.Latitude, result.Longitude)

	return result.Latitude, result.Longitude, result.Name, result.Country, nil
}

// getWeatherData получает данные о погоде от Open-Meteo Weather API
func (wb *WeatherBot) getWeatherData(lat, lon float64) (*WeatherResponse, error) {
	// Формируем URL для запроса к Weather API
	baseURL := "https://api.open-meteo.com/v1/forecast"
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%.4f", lat))
	params.Add("longitude", fmt.Sprintf("%.4f", lon))
	params.Add("current", "temperature_2m,apparent_temperature,relative_humidity_2m,weather_code,surface_pressure,wind_speed_10m,wind_direction_10m")
	params.Add("timezone", "auto")

	fullURL := baseURL + "?" + params.Encode()

	// Делаем HTTP запрос
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make weather request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	// Декодируем JSON ответ
	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	log.Printf("Weather data retrieved from Open-Meteo API\n")
	return &weather, nil
}

// getWeatherDescription возвращает описание погоды по WMO коду
func getWeatherDescription(code int) (string, string) {
	switch code {
	case 0:
		return "☀️", "Ясно"
	case 1:
		return "🌤️", "Преимущественно ясно"
	case 2:
		return "⛅", "Переменная облачность"
	case 3:
		return "☁️", "Пасмурно"
	case 45, 48:
		return "🌫️", "Туман"
	case 51, 53, 55:
		return "🌦️", "Морось"
	case 56, 57:
		return "🌨️", "Ледяная морось"
	case 61, 63, 65:
		return "🌧️", "Дождь"
	case 66, 67:
		return "🌨️", "Ледяной дождь"
	case 71, 73, 75:
		return "❄️", "Снег"
	case 77:
		return "❄️", "Снежные зерна"
	case 80, 81, 82:
		return "🌦️", "Ливень"
	case 85, 86:
		return "🌨️", "Снегопад"
	case 95:
		return "⛈️", "Гроза"
	case 96, 99:
		return "⛈️", "Гроза с градом"
	default:
		return "🌤️", "Неизвестно"
	}
}

// formatWeatherMessage форматирует данные о погоде в красивое сообщение
func (wb *WeatherBot) formatWeatherMessage(weather *WeatherResponse, cityName, country string) string {
	// Получаем описание погоды по коду
	weatherEmoji, weatherDesc := getWeatherDescription(weather.Current.WeatherCode)

	// Форматируем направление ветра
	windDirection := getWindDirection(weather.Current.WindDirection10m)

	return fmt.Sprintf(`%s Погода в %s, %s

🌡️ Температура: %.1f°C (ощущается как %.1f°C)
📝 Описание: %s
💨 Ветер: %.1f м/с (%s)
💧 Влажность: %d%%
📊 Давление: %.0f гПа

🕐 Обновлено: %s
📡 Данные: Open-Meteo.com`,
		weatherEmoji,
		cityName,
		country,
		weather.Current.Temperature2m,
		weather.Current.ApparentTemperature,
		weatherDesc,
		weather.Current.WindSpeed10m,
		windDirection,
		weather.Current.RelativeHumidity2m,
		weather.Current.SurfacePressure,
		time.Now().Format("15:04"),
	)
}

// getWindDirection возвращает направление ветра по градусам
func getWindDirection(degrees float64) string {
	directions := []string{"С", "ССВ", "СВ", "ВСВ", "В", "ВЮВ", "ЮВ", "ЮЮВ", "Ю", "ЮЮЗ", "ЮЗ", "ЗЮЗ", "З", "ЗСЗ", "СЗ", "ССЗ"}
	index := int((degrees+11.25)/22.5) % 16
	return directions[index]
}

// defaultHandler обрабатывает все остальные сообщения
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: `❓ Не понимаю эту команду.

Используйте:
• /weather <город> - узнать погоду
• /popular - популярные города
• /start - показать помощь`,
	})
	if err != nil {
		log.Printf("Failed to send default message: %v\n", err)
		return
	}
}
