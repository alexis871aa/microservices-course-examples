package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/olezhek28/microservices-course-examples/week_8/http/pkg/models"
)

const (
	httpPort     = "8080"
	urlParamCity = "city"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	// Создаем хранилище для данных о погоде
	storage := models.NewWeatherStorage()

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Health endpoint для проверки работоспособности
	r.Get("/api/v1/weather/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"ok","service":"weather-api"}`))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	// Определяем маршруты
	r.Route("/api/v1/weather", func(r chi.Router) {
		r.Get("/{city}", getWeatherHandler(storage))
		r.Put("/{city}", updateWeatherHandler(storage))
	})

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", httpPort), // Слушаем на всех интерфейсах для контейнеризации
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}

// getWeatherHandler обрабатывает запросы на получение информации о погоде для города
func getWeatherHandler(storage *models.WeatherStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		city := chi.URLParam(r, urlParamCity)
		if city == "" {
			http.Error(w, "City parameter is required", http.StatusBadRequest)
			return
		}

		weather := storage.GetWeather(city)
		if weather == nil {
			http.Error(w, fmt.Sprintf("Weather for city '%s' not found", city), http.StatusNotFound)
			return
		}

		render.JSON(w, r, weather)
	}
}

// updateWeatherHandler обрабатывает запросы на обновление информации о погоде для города
func updateWeatherHandler(storage *models.WeatherStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		city := chi.URLParam(r, urlParamCity)
		if city == "" {
			http.Error(w, "City parameter is required", http.StatusBadRequest)
			return
		}

		// Декодируем данные из тела запроса
		var weatherUpdate models.Weather
		if err := json.NewDecoder(r.Body).Decode(&weatherUpdate); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Устанавливаем имя города из URL-параметра
		weatherUpdate.City = city

		// Устанавливаем время обновления
		weatherUpdate.UpdatedAt = time.Now()

		// Обновляем информацию о погоде
		storage.UpdateWeather(&weatherUpdate)

		// Возвращаем обновленные данные
		render.JSON(w, r, weatherUpdate)
	}
}
