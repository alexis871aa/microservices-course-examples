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

	"github.com/alexis871aa/microservices-course-examples/week_1/http_chi/pkg/models"
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
	r.Use(middleware.Logger)                             // middleware, которая логгирует старт и конец каждого запроса, полезна, когда хотим логгировать каждый вызов
	r.Use(middleware.Recoverer)                          // middleware, при вызове наших ручек ,перехватывает панику, то есть по сути спасает от крит ошибок
	r.Use(middleware.Timeout(10 * time.Second))          // middleware, которая позволяет задать таймаут, то есть время ожидания, обрубает ответ, если ручка не смогла справится с запросом, через заданное количество времени
	r.Use(render.SetContentType(render.ContentTypeJSON)) // middleware, позволяет задать заголовки в http

	// Определяем маршруты, регистрируем роуты, то есть по сути те самые ручки
	r.Route("/api/v1/weather", func(r chi.Router) { // базовый url
		r.Get("/{city}", getWeatherHandler(storage))    // ручка для получения метод get, в path параметрах {city}
		r.Put("/{city}", updateWeatherHandler(storage)) // ручка для изменения метод put, в path парамтерах также ожидаем {city}
	})

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,                 // подсовываем нашему серверу главный обработчик r, который мы до этого сконфигурировали
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине, так как server.ListenAndServe() это блокирующий вызов
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown, это концепция акккуратного завершения программы и реакции на какие-то действия операционной системы
	quit := make(chan os.Signal, 1) // создаём канал типа os.Signal
	// нам интересны несколько:
	// SIGKILL - это сигнал, который говорит приложению, я даже тебя спрашивать не буду, а просто грохну. Вот этот сигнал перехватить нельзя (т.н. принудительное завершение работы)
	// SIGTERM - это сигнал, когда операционка приходит и говорит там такая ситуация - мы тебя хотим закрыть, вот такой сигнал лучше перехватывать и инициировать аккуратное завершение работы
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // тут мы висим на чтении и можем прочитать эти сигналы, как только что-то происходит выполнение кода идёт дальше

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx) // у самого сервера вызываем метод Shutdown и передаём туда контекст, таким образом, мы новые запросы отменяем, т.к. закрываемся, а старые запросы пытаемся завершить, но в пределах таймаута
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

		// render - это пакет в chi, который сериализует данные
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
