// Package logger предоставляет dual-write логгер с использованием zapcore.Tee архитектуры
//
// АРХИТЕКТУРА ЛОГГЕРА:
//
// Логгер использует zapcore.NewTee для параллельной записи в два назначения:
// 1. Stdout (для Kubernetes/контейнерных окружений)
// 2. OpenTelemetry коллектор (для централизованного сбора логов)
//
// ПОТОК ДАННЫХ:
//
//		Application
//		    ↓ (logger.Info/Error)
//		zap.Logger
//		    ↓
//		zapcore.Tee
//		   ↙        ↘
//	 StdoutCore   SimpleOTLPCore
//		   ↓             ↓
//	 os.Stdout   SimpleOTLPWriter
//		               ↓
//		        zapcore.BufferedWriteSyncer
//		               ↓
//		         OTLP Collector (gRPC)
//
// КОМПОНЕНТЫ:
//
// 1. StdoutCore - стандартный zap core для вывода в консоль
// 2. SimpleOTLPCore - преобразует zap Entry в OpenTelemetry Record
// 3. SimpleOTLPWriter - отправляет OTLP Records в коллектор
// 4. BufferedWriteSyncer - буферизация для асинхронной отправки
//
// ОСОБЕННОСТИ:
//
// - Graceful degradation: при недоступности OTLP коллектора stdout продолжает работать
// - Метрики: отслеживание sent/dropped записей для мониторинга
// - Батчирование: OTLP SDK автоматически группирует записи для эффективной отправки
// - Таймауты: 500ms лимит для предотвращения блокировки приложения
package logger

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelLog "go.opentelemetry.io/otel/log"
	otelLogSdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Глобальные переменные пакета
var (
	global       *zap.Logger                // глобальный экземпляр логгера
	initOnce     sync.Once                  // обеспечивает единократную инициализацию
	level        zap.AtomicLevel            // уровень логирования (может изменяться динамически)
	otelProvider *otelLogSdk.LoggerProvider // OTLP provider для graceful shutdown
)

// Константы конфигурации OTLP
const (
	otlpEndpoint       = "localhost:4317" // адрес OTLP коллектора
	serviceName        = "note-service"   // имя сервиса в телеметрии
	serviceEnvironment = "dev"            // окружение для фильтрации логов
)

// Таймауты
const (
	shutdownTimeout = 2 * time.Second // таймаут для graceful shutdown OTLP provider
)

// Init инициализирует глобальный логгер с Tee архитектурой.
// Поддерживает одновременную запись в stdout и OTLP коллектор.
//
// Параметры:
//   - logLevel: уровень логирования ("debug", "info", "warn", "error")
//   - asJSON: формат вывода (true - JSON, false - консольный)
//   - enableOTLP: включение отправки в OpenTelemetry коллектор
func Init(logLevel string, asJSON bool, enableOTLP bool) error {
	initOnce.Do(func() {
		level = zap.NewAtomicLevelAt(parseLevel(logLevel))
		cores := buildCores(asJSON, enableOTLP)
		global = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
	})

	if global == nil {
		return fmt.Errorf("logger init failed")
	}

	return nil
}

// buildCores создает слайс cores для zapcore.Tee.
// Всегда включает stdout core, опционально добавляет OTLP core.
func buildCores(asJSON bool, enableOTLP bool) []zapcore.Core {
	cores := []zapcore.Core{
		createStdoutCore(asJSON),
	}

	if enableOTLP {
		if otlpCore := createOTLPCore(); otlpCore != nil {
			cores = append(cores, otlpCore)
		}
	}

	return cores
}

// createStdoutCore создает core для записи в stdout/stderr.
// Поддерживает JSON и консольный формат вывода.
func createStdoutCore(asJSON bool) zapcore.Core {
	config := buildEncoderConfig()
	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
}

// createOTLPCore создает core для отправки в OpenTelemetry коллектор.
// При ошибке подключения возвращает nil (graceful degradation).
func createOTLPCore() *SimpleOTLPCore {
	otlpLogger, err := createOTLPLogger(otlpEndpoint)
	if err != nil {
		// Логирование ошибки невозможно, так как логгер еще не инициализирован
		return nil
	}

	// Прямо передаём OTLP-логгер в core. Буферизацию делает OTLP SDK (BatchProcessor).
	return NewSimpleOTLPCore(otlpLogger, level)
}

// createOTLPLogger создает OTLP логгер с настроенным экспортером и ресурсами.
// Использует BatchProcessor для эффективной отправки логов.
func createOTLPLogger(endpoint string) (otelLog.Logger, error) {
	ctx := context.Background()

	exporter, err := createOTLPExporter(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	rs, err := createResource(ctx)
	if err != nil {
		return nil, err
	}

	provider := otelLogSdk.NewLoggerProvider(
		otelLogSdk.WithResource(rs),
		otelLogSdk.WithProcessor(otelLogSdk.NewBatchProcessor(exporter)),
	)
	otelProvider = provider // сохраняем для shutdown

	return provider.Logger("app"), nil
}

// createOTLPExporter создает gRPC экспортер для OTLP коллектора
func createOTLPExporter(ctx context.Context, endpoint string) (*otlploggrpc.Exporter, error) {
	return otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(), // для разработки, в продакшене следует использовать TLS
	)
}

// createResource создает метаданные сервиса для телеметрии
func createResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("deployment.environment", serviceEnvironment),
		),
	)
}

// buildEncoderConfig настраивает формат вывода логов с нужными полями
func buildEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
}

// Info записывает лог уровня INFO.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Info(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Info(msg, fields...)
	}
}

// Error записывает лог уровня ERROR.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Error(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Error(msg, fields...)
	}
}

// Sync принудительно сбрасывает все буферизованные логи.
// Вызывает sync для всех cores (stdout + OTLP).
func Sync() error {
	if global != nil {
		return global.Sync()
	}

	return nil
}

// Close корректно завершает работу логгера.
// Останавливает OTLP provider с таймаутом для отправки оставшихся логов.
func Close() error {
	if otelProvider != nil {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = otelProvider.Shutdown(ctx)
	}

	return nil
}

// parseLevel преобразует строковое значение в zapcore.Level
func parseLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
