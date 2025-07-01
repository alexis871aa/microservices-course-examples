package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestsTotal общий счетчик запросов по методам и статусам
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ufo_requests_total",
			Help: "Total number of UFO service requests",
		},
		[]string{"method", "status"},
	)

	// SightingsTotal счетчик созданных наблюдений
	SightingsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "ufo_sightings_total",
			Help: "Total number of UFO sightings created",
		},
	)

	// AnalysisRequestsTotal счетчик запросов на анализ
	AnalysisRequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "ufo_analysis_requests_total",
			Help: "Total number of UFO analysis requests",
		},
	)

	// AnalysisDuration гистограмма времени выполнения анализа
	AnalysisDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "ufo_analysis_duration_seconds",
			Help:    "Duration of UFO analysis operations",
			Buckets: prometheus.DefBuckets,
		},
	)

	// ActiveConnections gauge активных соединений
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ufo_active_connections",
			Help: "Number of active gRPC connections",
		},
	)

	// RequestDuration гистограмма времени выполнения запросов
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ufo_request_duration_seconds",
			Help:    "Duration of gRPC requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

// RegisterMetrics регистрирует все метрики UFO сервиса
func RegisterMetrics() {
	// Метрики уже зарегистрированы через promauto
	// Эта функция нужна для явной инициализации
}
