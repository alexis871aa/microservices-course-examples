package metrics

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	defaultTimeout = 5 * time.Second
)

var (
	exporter      *otlpmetricgrpc.Exporter
	meterProvider *metric.MeterProvider
)

type Config interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
}

// InitProvider инициализирует глобальный провайдер метрик OpenTelemetry
func InitProvider(ctx context.Context, cfg Config) error {
	var err error

	// Создаем экспортер для отправки метрик в OTLP коллектор
	exporter, err = otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint()),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithTimeout(defaultTimeout),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create metrics exporter")
	}

	// Создаем провайдер метрик
	meterProvider = metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(
				exporter,
				metric.WithInterval(cfg.CollectorInterval()),
			),
		),
	)

	// Устанавливаем глобальный провайдер метрик
	otel.SetMeterProvider(meterProvider)

	return nil
}

// GetMeterProvider возвращает текущий провайдер метрик
func GetMeterProvider() *metric.MeterProvider {
	return meterProvider
}

// Shutdown закрывает провайдер метрик и экспортер
func Shutdown(ctx context.Context) error {
	if meterProvider == nil && exporter == nil {
		return nil
	}

	var err error

	if meterProvider != nil {
		err = meterProvider.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to shutdown meter provider")
		}
	}

	if exporter != nil {
		err = exporter.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to shutdown exporter")
		}
	}

	return nil
}
