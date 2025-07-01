package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/olezhek28/microservices-course-examples/week_1/http_chi/pkg/models"
)

const (
	serverURL         = "http://localhost:8080"
	weatherAPIPath    = "/api/v1/weather/%s"
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
	requestTimeout    = 5 * time.Second
	defaultCityName   = "Moscow"
	defaultMinTemp    = -10
	defaultMaxTemp    = 40
)

// httpClient - HTTP клиент с таймаутом
var httpClient = &http.Client{
	Timeout: requestTimeout,
}

// getWeather получает информацию о погоде для указанного города
func getWeather(ctx context.Context, city string) (*models.Weather, error) {
	// Создаем запрос с контекстом
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s"+weatherAPIPath, serverURL, city),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("создание GET-запроса: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("выполнение GET-запроса: %w", err)
	}
	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			log.Printf("ошибка закрытия тела ответа: %v\n", cerr)
			return
		}
	}()

	// Если статус ответа 404, считаем что данных нет
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	// Читаем тело ответа для любого ответа (для логирования при ошибке)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("чтение тела ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("получение данных о погоде (статус %d): %s", resp.StatusCode, string(body))
	}

	var weather models.Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, fmt.Errorf("декодирование JSON: %w", err)
	}

	return &weather, nil
}

// updateWeather обновляет данные о погоде для указанного города
func updateWeather(ctx context.Context, city string, weather *models.Weather) (*models.Weather, error) {
	// Кодируем данные о погоде в JSON
	jsonData, err := json.Marshal(weather)
	if err != nil {
		return nil, fmt.Errorf("кодирование JSON: %w", err)
	}

	// Создаем PUT-запрос с контекстом
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("%s"+weatherAPIPath, serverURL, city),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("создание PUT-запроса: %w", err)
	}
	req.Header.Set(contentTypeHeader, contentTypeJSON)

	// Выполняем запрос
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("выполнение PUT-запроса: %w", err)
	}
	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			log.Printf("ошибка закрытия тела ответа: %v\n", cerr)
			return
		}
	}()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("чтение тела ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("обновление данных о погоде (статус %d): %s", resp.StatusCode, string(body))
	}

	// Декодируем ответ
	var updatedWeather models.Weather
	err = json.Unmarshal(body, &updatedWeather)
	if err != nil {
		return nil, fmt.Errorf("декодирование JSON: %w", err)
	}

	return &updatedWeather, nil
}

// generateRandomWeather создает случайные данные о погоде
func generateRandomWeather() *models.Weather {
	return &models.Weather{
		Temperature: gofakeit.Float64Range(defaultMinTemp, defaultMaxTemp),
	}
}

func main() {
	ctx := context.Background()

	log.Println("=== Тестирование API для работы с данными о погоде ===")
	log.Println()

	// 1. Пытаемся получить данные о погоде (которых пока нет)
	log.Printf("🌦️ Получение данных о погоде для города %s\n", defaultCityName)
	log.Println("===================================================")

	weather, err := getWeather(ctx, defaultCityName)
	if err != nil {
		log.Printf("❌ Ошибка: %v\n", err)
		return
	}

	log.Printf("Данные о погоде для города %s: %+v\n", defaultCityName, weather)

	// 2. Обновляем данные о погоде
	log.Printf("🔄 Обновление данных о погоде для города %s\n", defaultCityName)
	log.Println("=====================================================")

	newWeather := generateRandomWeather()

	updatedWeather, err := updateWeather(ctx, defaultCityName, newWeather)
	if err != nil {
		log.Printf("❌ Ошибка при обновлении погоды: %v\n", err)
		return
	}
	log.Printf("✅ Данные о погоде обновлены: %+v\n", updatedWeather)

	// 3. Получаем обновленные данные о погоде
	log.Printf("🌦️ Получение обновленных данных о погоде для города %s\n", defaultCityName)
	log.Println("===========================================================")

	weather, err = getWeather(ctx, defaultCityName)
	if err != nil {
		log.Printf("❌ Ошибка при получении погоды: %v\n", err)
		return
	}

	if weather == nil {
		log.Printf("❌ Неожиданно: данные о погоде отсутствуют после обновления\n")
		return
	}

	log.Printf("✅ Получены данные о погоде: %+v\n", weather)
	log.Println("Тестирование завершено успешно!")
}
