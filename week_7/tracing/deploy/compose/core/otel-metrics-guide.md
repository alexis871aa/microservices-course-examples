# Руководство по системе метрик с OpenTelemetry

Данное руководство объясняет как устроен сбор метрик в микросервисной архитектуре с использованием OpenTelemetry, начиная с основных понятий и заканчивая подробным описанием ключевых компонентов.

## 1. Общая архитектура системы мониторинга

Наша система мониторинга использует современный стек инструментов, основанный на OpenTelemetry:

```
┌───────────────┐         ┌──────────────────┐         ┌───────────────┐         ┌──────────────┐
│  Микросервисы │  ──►    │  OpenTelemetry   │  ──►    │  Prometheus   │  ──►    │   Grafana    │
│ (инструменты) │   OTLP   │    Collector     │  Remote │  (хранилище   │   Query │ (визуализация│ 
└───────────────┘         └──────────────────┘  Write  └───────────────┘         └──────────────┘
```

Основные компоненты:

1. **Микросервисы** (Order, Assembly, Payment и др.) собирают телеметрию с помощью OpenTelemetry SDK
2. **OpenTelemetry Collector** получает данные, обрабатывает их и передает в системы хранения
3. **Prometheus** хранит метрики и обеспечивает их запрос
4. **Grafana** визуализирует метрики на дашбордах

## 2. Основные понятия OpenTelemetry

### 2.1 Что такое OpenTelemetry

OpenTelemetry (также обозначается как OTel) - это набор API, библиотек, агентов и инструментария для сбора и экспорта телеметрии из программного обеспечения. Телеметрия делится на три основных типа:

- **Метрики** (Metrics) - числовые данные для измерения производительности
- **Трассировка** (Traces) - информация о пути выполнения запросов через систему
- **Логи** (Logs) - записи о событиях в системе

В нашем проекте мы в первую очередь используем OpenTelemetry для сбора метрик.

### 2.2 Типы метрик в OpenTelemetry

В OpenTelemetry существуют следующие типы метрик:

- **Counter** (Счетчик) - значение, которое может только увеличиваться. Примеры: количество заказов, число HTTP-запросов.
- **Gauge** (Датчик) - значение, которое может увеличиваться и уменьшаться. Примеры: использование памяти, количество активных соединений.
- **Histogram** (Гистограмма) - распределение значений. Пример: время ответа сервера.
- **UpDownCounter** - счетчик, который может как увеличиваться, так и уменьшаться.

### 2.3 Протокол OTLP

OpenTelemetry Protocol (OTLP) - это стандартизированный протокол для передачи телеметрии. Он поддерживает передачу по gRPC (бинарный формат) и HTTP/protobuf.

В нашем проекте мы используем OTLP для отправки метрик из сервисов напрямую в OpenTelemetry Collector, минуя промежуточные этапы, что повышает производительность и надежность системы.

## 3. Компоненты системы сбора метрик

### 3.1 OpenTelemetry Collector

OpenTelemetry Collector - это промежуточное звено между источниками телеметрии и системами хранения. Его основные функции:

- **Получение данных** из различных источников
- **Обработка данных** (фильтрация, агрегация, обогащение)
- **Экспорт данных** в различные системы хранения

#### 3.1.1 Структура конфигурации Collector

Конфигурация коллектора находится в файле `deploy/compose/core/otel/otel-collector-config.yaml` и состоит из нескольких основных секций:

- **receivers** - получатели данных (откуда брать телеметрию)
- **processors** - обработчики данных (как их трансформировать)
- **exporters** - отправители данных (куда их передавать)
- **pipelines** - определение потоков данных от получателей через обработчики к отправителям

#### 3.1.2 Receivers (Получатели)

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
```

Эта секция настраивает OpenTelemetry Collector для приема данных по протоколу OTLP:
- По gRPC на порту 4317
- По HTTP на порту 4318

#### 3.1.3 Processors (Обработчики)

```yaml
processors:
  batch:
    send_batch_size: 10000
    timeout: 10s
  memory_limiter:
    check_interval: 1s
    limit_percentage: 75
    spike_limit_percentage: 15
  resourcedetection:
    detectors: [env, system]
    timeout: 5s
```

Эта секция определяет обработчики для данных:
- **batch** - группирует данные для эффективной отправки
- **memory_limiter** - предотвращает перегрузку памяти
- **resourcedetection** - добавляет метаданные о системе к телеметрии

#### 3.1.4 Exporters (Экспортеры)

```yaml
exporters:
  prometheusremotewrite:
    endpoint: http://prometheus:9090/api/v1/write
    tls:
      insecure: true
  logging:
    loglevel: debug
```

Эта секция определяет куда отправлять метрики:
- **prometheusremotewrite** - отправляет метрики напрямую в Prometheus по Remote Write API
- **logging** - выводит данные в лог (для отладки)

#### 3.1.5 Service и Pipelines

```yaml
service:
  extensions: [health_check, pprof, zpages]
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch, resourcedetection]
      exporters: [prometheusremotewrite, logging]
```

Эта секция связывает все компоненты в единую систему:
- Порядок обработки для метрик: получение через OTLP → обработка → отправка в Prometheus

### 3.2 Prometheus

Prometheus - это система мониторинга и база данных временных рядов, которая хранит метрики и обеспечивает их запрос.

#### 3.2.1 Структура конфигурации Prometheus

Конфигурация Prometheus находится в файле `deploy/compose/core/prometheus/prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets: ["localhost:9090"]

  - job_name: "otel-collector"
    static_configs:
      - targets: ["otel-collector:8888"]
```

Ключевые моменты:
- **scrape_interval** - как часто собирать метрики (15 секунд)
- **scrape_configs** - определяет откуда собирать метрики
  - Prometheus собирает метрики о себе самом
  - Также собираются метрики о состоянии OpenTelemetry Collector

### 3.3 Grafana

Grafana - это платформа для визуализации и анализа метрик. Она позволяет создавать дашборды с графиками, таблицами и алертами.

#### 3.3.1 Настройка источников данных

В файле `deploy/compose/core/grafana/provisioning/datasources/prometheus.yml` настраивается подключение к Prometheus:

```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
    version: 1
```

#### 3.3.2 Настройка дашбордов

Дашборды автоматически загружаются из директории `/var/lib/grafana/dashboards` согласно настройкам в `deploy/compose/core/grafana/provisioning/dashboards/dashboard.yml`:

```yaml
apiVersion: 1
providers:
  - name: 'Microservices Dashboard'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    editable: true
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /var/lib/grafana/dashboards
```

В директории `deploy/compose/core/grafana/dashboards` находятся JSON-файлы с описанием дашбордов.

## 4. Инструментирование микросервисов

### 4.1 Структура библиотек метрик в микросервисах

В каждом микросервисе (Order, Assembly, Payment, Inventory) есть пакет `metric`, который отвечает за сбор и отправку метрик.

#### 4.1.1 Инициализация метрик

Типичный код инициализации в файле `metric/metric.go`:

```go
func InitMetrics(ctx context.Context, cfg config.MetricServerConfig) error {
    // Создание экспортера OTLP
    exporter, err := otlpmetricgrpc.New(ctx,
        otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint()),
        otlpmetricgrpc.WithInsecure(),
    )
    
    // Настройка провайдера метрик
    provider := sdkmetric.NewMeterProvider(
        sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter,
            sdkmetric.WithInterval(cfg.CollectorInterval()))),
    )
    
    // Установка глобального провайдера
    otel.SetMeterProvider(provider)
    
    // Создание метрик
    // ...
}
```

#### 4.1.2 Определение метрик

Каждый сервис определяет свои метрики в зависимости от бизнес-логики:

```go
// Пример метрик из сервиса Order
orderCounter, err = meter.Int64Counter(
    metricOrdersTotal,
    metric.WithDescription("Количество созданных заказов"),
)

// Пример метрик из сервиса Assembly
assemblyDuration, err = meter.Float64Histogram(
    metricAssemblyDuration,
    metric.WithDescription("Время сборки заказов в секундах"),
)
```

#### 4.1.3 Методы регистрации метрик

В каждом сервисе определены методы для регистрации соответствующих метрик:

```go
// В сервисе Order
func CountOrderCreated(ctx context.Context, status string) {
    orderCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("status", status)))
}

// В сервисе Assembly
func ObserveAssemblyDuration(ctx context.Context, duration time.Duration) {
    assemblyDuration.Record(ctx, duration.Seconds())
}
```

### 4.2 Интеграция в бизнес-логику

Метрики интегрируются в ключевые точки бизнес-логики для измерения важных параметров:

```go
// Пример из обработчика в Assembly
func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
    // ...обработка...
    
    startTime := time.Now()
    // ...бизнес-логика...
    assemblyDuration := time.Since(startTime)
    
    // Записываем метрики сборки
    metric.ObserveAssemblyDuration(ctx, assemblyDuration)
    metric.CountAssemblyOrder(ctx)
    
    // ...завершение обработки...
}
```

### 4.3 Middleware и интерцепторы

Для автоматического сбора метрик HTTP-запросов используются middleware:

```go
// Middleware для сбора HTTP метрик
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := newResponseWriter(w)
        next.ServeHTTP(rw, r)
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(rw.statusCode)
        
        metric.ObserveHttpDuration(r.Context(), r.Method, r.URL.Path, status, duration)
    })
}
```

## 5. Конфигурация для OpenTelemetry

### 5.1 Структура конфигурации метрик

Конфигурация для метрик в каждом сервисе определяется через интерфейс `MetricServerConfig`:

```go
type MetricServerConfig interface {
    CollectorEndpoint() string     // Адрес коллектора OpenTelemetry
    CollectorInterval() time.Duration  // Интервал отправки метрик
}
```

Реализация обычно находится в файлах:
- `internal/config/env/metric.go`

### 5.2 Переменные окружения

Для настройки параметров метрик используются переменные окружения:

```
COLLECTOR_ENDPOINT=otel-collector:4317  # Адрес OpenTelemetry Collector
COLLECTOR_INTERVAL=10s                 # Интервал отправки метрик (10 секунд)
```

Эти переменные настраиваются в файлах `.env` или передаются через Docker Compose.

## 6. Docker Compose настройки

Конфигурация Docker Compose в файле `deploy/compose/core/docker-compose.yml` определяет все компоненты системы метрик:

```yaml
services:
  otel-collector:
    image: otel/opentelemetry-collector-contrib:0.100.0
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./otel/otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      - "${OTEL_GRPC_PORT:-4317}:4317"  # Порт для OTLP gRPC
      - "${OTEL_HTTP_PORT:-4318}:4318"  # Порт для OTLP HTTP
      # ...другие порты...

  prometheus:
    image: prom/prometheus:v2.49.1
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    # ...другие настройки...

  grafana:
    image: grafana/grafana:10.3.1
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/provisioning:/etc/grafana/provisioning
      - ./grafana/dashboards:/var/lib/grafana/dashboards
    # ...другие настройки...
```

## 7. Руководство по отладке метрик

### 7.1 Проверка отправки метрик

1. **Просмотр логов OpenTelemetry Collector**:
   ```
   docker logs otel-collector
   ```
   Коллектор в режиме отладки показывает все полученные метрики.

2. **Проверка метрик через Prometheus UI**:
   - Откройте `http://localhost:9090/graph`
   - Введите имя метрики в поле поиска (например, `order_service_orders_total`)

3. **Проверка через Grafana**:
   - Откройте `http://localhost:3000`
   - Войдите (логин/пароль по умолчанию: admin/admin)
   - Перейдите к дашборду Microservices

### 7.2 Общие проблемы и их решения

1. **Метрики не отображаются в Prometheus**:
   - Проверьте подключение между сервисом и OpenTelemetry Collector
   - Убедитесь, что имена метрик соответствуют ожидаемым
   - Проверьте, работает ли экспортер prometheusremotewrite

2. **Проблемы с подключением к коллектору**:
   - Убедитесь, что сервисы могут достичь коллектор по сети
   - Проверьте настройку COLLECTOR_ENDPOINT

3. **Высокое потребление ресурсов**:
   - Настройте параметры batch processor для уменьшения нагрузки
   - Увеличьте интервал отправки метрик (COLLECTOR_INTERVAL)

## 8. Заключение

Система сбора метрик на основе OpenTelemetry предоставляет мощный инструмент для мониторинга микросервисной архитектуры. Она позволяет:

- Собирать унифицированные метрики со всех сервисов
- Гибко настраивать обработку и хранение метрик
- Визуализировать данные для анализа производительности и выявления проблем

Правильная настройка и использование этих инструментов позволяет не только реагировать на проблемы, но и предотвращать их.
