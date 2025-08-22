// Package metric предоставляет инструменты для сбора и отправки метрик приложения
// с использованием OpenTelemetry SDK
package metric

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	namespace   = "my_space"     // Пространство имен для группировки метрик
	appName     = "my_app"       // Имя приложения для включения в имена метрик
	serviceName = "note-service" // Имя сервиса для идентификации источника метрик
)

// Глобальные переменные для хранения инструментов метрик
var (
	// meter - основной объект для создания инструментов метрик
	// Meter в OpenTelemetry - это фабрика для создания инструментов (счетчиков, гистограмм и т.д.)
	// Связывается с конкретным сервисом и версией библиотеки для изоляции метрик
	meter = otel.Meter(serviceName)

	// requestCounter - счетчик входящих запросов
	// Counter - монотонно возрастающее значение (никогда не уменьшается)
	// Используется для подсчета событий: запросы, ошибки, завершенные операции
	requestCounter metric.Int64Counter

	// responseCounter - счетчик исходящих ответов с атрибутами (статус, метод)
	// Позволяет группировать метрики по различным измерениям
	responseCounter metric.Int64Counter

	// histogramResponseTime - гистограмма времени ответа
	// Histogram - распределение значений по корзинам (buckets)
	// Позволяет вычислять процентили, медиану, среднее время ответа
	histogramResponseTime metric.Float64Histogram
)

// Init инициализирует все инструменты метрик
// Должна вызываться после настройки глобального MeterProvider
func Init(_ context.Context) error {
	var err error

	// Создание счетчика запросов
	// Имя метрики следует соглашениям: namespace_protocol_app_metric_unit
	// WithDescription добавляет описание для документации и UI
	requestCounter, err = meter.Int64Counter(
		namespace+"_grpc_"+appName+"_requests_total",
		metric.WithDescription("Количество запросов к серверу"),
	)
	if err != nil {
		return err
	}

	// Создание счетчика ответов с поддержкой атрибутов
	// Атрибуты (labels в Prometheus) позволяют группировать метрики
	// Например: status="success", method="GetNote"
	responseCounter, err = meter.Int64Counter(
		namespace+"_grpc_"+appName+"_responses_total",
		metric.WithDescription("Количество ответов от сервера"),
	)
	if err != nil {
		return err
	}

	// Создание гистограммы времени ответа
	// WithUnit указывает единицы измерения (важно для правильной интерпретации)
	// WithExplicitBucketBoundaries определяет бакеты для распределения значений
	// Бакеты выбраны для покрытия диапазона от микросекунд до секунд
	histogramResponseTime, err = meter.Float64Histogram(
		namespace+"_grpc_"+appName+"_histogram_response_time_seconds",
		metric.WithDescription("Время ответа от сервера"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(
			// Бакеты от 0.1ms до ~3.3s с экспоненциальным ростом
			0.0001, 0.0002, 0.0004, 0.0008, 0.0016, 0.0032, 0.0064, 0.0128,
			0.0256, 0.0512, 0.1024, 0.2048, 0.4096, 0.8192, 1.6384, 3.2768,
		),
	)
	if err != nil {
		return err
	}

	return nil
}

// IncRequestCounter увеличивает счетчик входящих запросов на 1
// Вызывается при получении каждого нового запроса
func IncRequestCounter(ctx context.Context) {
	requestCounter.Add(ctx, 1)
}

// IncResponseCounter увеличивает счетчик ответов с указанными атрибутами
// status - статус ответа (success, error, timeout и т.д.)
// method - имя вызываемого метода gRPC
// Атрибуты позволяют создавать срезы метрик для анализа
func IncResponseCounter(ctx context.Context, status, method string) {
	responseCounter.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("status", status),
			attribute.String("method", method),
		),
	)
}

// HistogramResponseTimeObserve записывает время выполнения запроса в гистограмму
// status - статус выполнения операции
// time - время выполнения в секундах
// Гистограмма автоматически распределяет значения по корзинам для последующего анализа
func HistogramResponseTimeObserve(ctx context.Context, status string, time float64) {
	histogramResponseTime.Record(ctx, time,
		metric.WithAttributes(
			attribute.String("status", status),
		),
	)
}

// InitOpenTelemetry инициализирует OpenTelemetry provider и настраивает экспорт метрик
// Возвращает MeterProvider, который нужно закрыть при завершении приложения
//
// Основные компоненты OpenTelemetry:
// 1. Resource - метаданные о сервисе (имя, версия, окружение)
// 2. Exporter - компонент для отправки данных во внешние системы
// 3. Reader - компонент для чтения метрик из инструментов
// 4. MeterProvider - центральная точка управления метриками
func InitOpenTelemetry() (*sdkmetric.MeterProvider, error) {
	// 1. Создаем OTLP gRPC экспортер для отправки метрик в OpenTelemetry Collector
	// OTLP (OpenTelemetry Protocol) - стандартный протокол для передачи телеметрии
	// Collector действует как промежуточный слой между приложением и системами мониторинга
	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithEndpoint("localhost:4317"),                // Стандартный порт OTLP gRPC
		otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()), // Без TLS для локальной разработки
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// 2. Создаем ресурс с метаданными о сервисе
	// Resource - это набор атрибутов, описывающих источник телеметрии
	// Эти атрибуты добавляются ко всем метрикам автоматически
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),     // Обязательный атрибут
			attribute.String("service.version", "1.0.0"),      // Версия сервиса
			attribute.String("deployment.environment", "dev"), // Окружение (dev/staging/prod)
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 3. Создаем MeterProvider с периодическим reader
	// MeterProvider - основной компонент SDK, управляющий всеми метриками
	// PeriodicReader - читает метрики с заданным интервалом и отправляет через экспортер
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(5*time.Second), // Интервал сбора и отправки метрик
			),
		),
	)

	// 4. Устанавливаем глобальный MeterProvider
	// Это позволяет получать Meter через otel.Meter() из любого места в коде
	// Глобальное состояние упрощает использование, но может усложнить тестирование
	otel.SetMeterProvider(meterProvider)

	log.Println("OpenTelemetry initialized successfully")
	return meterProvider, nil
}
