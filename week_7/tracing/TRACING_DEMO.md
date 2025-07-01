# Демонстрация распределенного трейсинга

Этот проект демонстрирует распределенный трейсинг с использованием OpenTelemetry, Jaeger и цепочки микросервисов.

## Архитектура системы

```
Client → UFO Service → Analysis Service → Classification Service
  :50051      :50052           :50053
```

### Сервисы

1. **UFO Service** (порт 50051) - основной сервис для работы с наблюдениями НЛО
2. **Analysis Service** (порт 50052) - сервис анализа наблюдений
3. **Classification Service** (порт 50053) - сервис классификации объектов

### Инфраструктура трейсинга

- **OpenTelemetry Collector** (порт 4317/4318) - сбор и обработка трейсов
- **Jaeger** (порт 16686) - хранение и визуализация трейсов

## Детальное взаимодействие сервисов

### UFO Service (порт 50051)

**Основные ручки:**
- `Create(CreateRequest) → CreateResponse` - создание наблюдения НЛО
- `Get(GetRequest) → GetResponse` - получение наблюдения по UUID
- `Update(UpdateRequest) → Empty` - обновление наблюдения
- `Delete(DeleteRequest) → Empty` - мягкое удаление наблюдения
- `AnalyzeSighting(AnalyzeSightingRequest) → AnalyzeSightingResponse` - **анализ наблюдения (трейсинг)**

**Трейсинг в UFO Service:**
```go
// ufo/internal/api/ufo/v1/analyze.go
func (a *api) AnalyzeSighting(ctx context.Context, req *ufoV1.AnalyzeSightingRequest) (*ufoV1.AnalyzeSightingResponse, error) {
    // Создаем спан для вызова Analysis сервиса
    ctx, span := tracing.StartSpan(ctx, "ufo.call_analysis")
    defer span.End()

    // Подключение к Analysis сервису с трейсинг интерцептором
    conn, err := grpc.NewClient(
        "localhost:50052",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("analysis-service")),
    )
    // ... вызов Analysis Service
}
```

**Инициализация трейсинга:**
```go
// ufo/internal/app/app.go
func (a *App) initTracing(ctx context.Context) error {
    err := tracing.InitTracer(ctx, config.AppConfig().Tracing)
    if err != nil {
        return err
    }
    
    closer.AddNamed("tracing", func(ctx context.Context) error {
        return tracing.ShutdownTracer(ctx)
    })
    
    return nil
}

// gRPC сервер с server интерцептором
func (a *App) initGRPCServer(ctx context.Context) error {
    a.grpcServer = grpc.NewServer(
        grpc.Creds(insecure.NewCredentials()),
        grpc.UnaryInterceptor(tracing.UnaryServerInterceptor("ufo-service")),
    )
    // ...
}
```

### Analysis Service (порт 50052)

**Основные ручки:**
- `AnalyzeSighting(AnalyzeSightingRequest) → AnalyzeSightingResponse` - анализ наблюдения

**Логика взаимодействия:**
1. Получает запрос от UFO Service с UUID наблюдения
2. Вызывает UFO Service для получения данных наблюдения (`Get`)
3. Вызывает Classification Service для классификации объекта
4. Возвращает результат анализа

**Трейсинг в Analysis Service:**
```go
// analysis/internal/service/analysis.go
func (s *AnalysisService) AnalyzeSighting(ctx context.Context, uuid string) (string, string, float32, error) {
    // Спан для получения данных о наблюдении
    ctx, span := tracing.StartSpan(ctx, "analysis.get_sighting")
    
    // Получаем данные о наблюдении из UFO сервиса
    sighting, err := s.ufoClient.Get(ctx, &ufoV1.GetRequest{Uuid: uuid})
    if err != nil {
        span.End()
        return "", "", 0, fmt.Errorf("failed to get sighting: %w", err)
    }
    span.End()

    // Спан для классификации
    ctx, span = tracing.StartSpan(ctx, "analysis.classify")
    defer span.End()

    // Отправляем данные в Classification сервис
    classification, err := s.classificationClient.ClassifyObject(ctx, &classificationV1.ClassifyObjectRequest{
        Description:     sighting.Sighting.Info.Description,
        Color:           sighting.Sighting.Info.Color.GetValue(),
        DurationSeconds: sighting.Sighting.Info.DurationSeconds.GetValue(),
    })
    // ...
}
```

**Инициализация с клиентскими интерцепторами:**
```go
// analysis/internal/service/analysis.go
func NewAnalysisService() (*AnalysisService, error) {
    // Подключение к UFO сервису
    ufoConn, err := grpc.NewClient(
        "localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("ufo-service")),
    )

    // Подключение к Classification сервису
    classificationConn, err := grpc.NewClient(
        "localhost:50053",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("classification-service")),
    )
    // ...
}
```

### Classification Service (порт 50053)

**Основные ручки:**
- `ClassifyObject(ClassifyObjectRequest) → ClassifyObjectResponse` - классификация объекта

**Логика классификации:**
```go
// classification/internal/service/classification.go
func (s *ClassificationService) ClassifyObject(ctx context.Context, req *classificationV1.ClassifyObjectRequest) (*classificationV1.ClassifyObjectResponse, error) {
    // Создаем спан для внутренней логики
    ctx, span := tracing.StartSpan(ctx, "classification.analyze")
    defer span.End()

    // Логика классификации по описанию
    objectType := "unknown"
    confidence := float32(0.5)
    explanation := "Базовая классификация"

    description := strings.ToLower(req.Description)
    
    // Классификация по форме
    if strings.Contains(description, "треугольн") {
        objectType = "triangular_craft"
        confidence = 0.8
        explanation = "Треугольная форма объекта"
    } else if strings.Contains(description, "диск") || strings.Contains(description, "тарелка") {
        objectType = "classic_saucer"
        confidence = 0.9
        explanation = "Классическая форма летающей тарелки"
    } else if strings.Contains(description, "сфер") || strings.Contains(description, "шар") {
        objectType = "orb"
        confidence = 0.7
        explanation = "Сферическая форма объекта"
    }

    // Корректировка уверенности по цвету и длительности
    // ...
}
```

## Реализация OpenTelemetry интерцепторов

### Платформенный пакет трейсинга

Все интерцепторы реализованы в `platform/pkg/tracing/`:

**Инициализация трейсера:**
```go
// platform/pkg/tracing/tracer.go
func InitTracer(ctx context.Context, cfg Config) error {
    serviceName = cfg.ServiceName()

    // Создаем экспортер для отправки трейсов в OpenTelemetry Collector
    exporter, err := otlptracegrpc.New(
        ctx,
        otlptracegrpc.WithEndpoint(cfg.CollectorEndpoint()),
        otlptracegrpc.WithInsecure(),
        // ... другие настройки
    )

    // Создаем ресурс с метаданными сервиса
    attributeResource, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName(cfg.ServiceName()),
            semconv.ServiceVersion(cfg.ServiceVersion()),
            attribute.String("environment", cfg.Environment()),
        ),
        // ...
    )

    // Создаем провайдер трейсов
    tracerProvider := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(attributeResource),
    )

    otel.SetTracerProvider(tracerProvider)
    otel.SetTextMapPropagator(propagation.TraceContext{})
    
    return nil
}
```

**Server интерцептор:**
```go
// platform/pkg/tracing/interceptors.go
func UnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // Извлекаем контекст трейсинга из входящего запроса
        ctx = otel.GetTextMapPropagator().Extract(ctx, &metadataCarrier{md: getIncomingMetadata(ctx)})
        
        // Создаем спан для обработки запроса
        ctx, span := otel.Tracer(serviceName).Start(ctx, info.FullMethod)
        defer span.End()

        // Вызываем обработчик
        resp, err := handler(ctx, req)
        if err != nil {
            span.RecordError(err)
            span.SetStatus(codes.Error, err.Error())
        }

        return resp, err
    }
}
```

**Client интерцептор:**
```go
func UnaryClientInterceptor(targetService string) grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        // Создаем спан для исходящего запроса
        ctx, span := otel.Tracer(serviceName).Start(ctx, method)
        defer span.End()

        // Добавляем атрибуты
        span.SetAttributes(
            attribute.String("rpc.service", targetService),
            attribute.String("rpc.method", method),
        )

        // Внедряем контекст трейсинга в метаданные
        md := make(metadata.MD)
        otel.GetTextMapPropagator().Inject(ctx, &metadataCarrier{md: md})
        ctx = metadata.NewOutgoingContext(ctx, md)

        // Выполняем вызов
        err := invoker(ctx, method, req, reply, cc, opts...)
        if err != nil {
            span.RecordError(err)
            span.SetStatus(codes.Error, err.Error())
        }

        return err
    }
}
```

**Создание спанов:**
```go
func StartSpan(ctx context.Context, operationName string) (context.Context, trace.Span) {
    return otel.Tracer(serviceName).Start(ctx, operationName)
}
```

### Пропагация контекста

**Metadata Carrier для gRPC:**
```go
type metadataCarrier struct {
    md metadata.MD
}

func (mc *metadataCarrier) Get(key string) string {
    values := mc.md.Get(key)
    if len(values) == 0 {
        return ""
    }
    return values[0]
}

func (mc *metadataCarrier) Set(key, value string) {
    mc.md.Set(key, value)
}

func (mc *metadataCarrier) Keys() []string {
    keys := make([]string, 0, len(mc.md))
    for k := range mc.md {
        keys = append(keys, k)
    }
    return keys
}
```

## Цепочка трейсов

При вызове `AnalyzeSighting` создается следующая цепочка спанов:

1. **Client → UFO Service**: gRPC запрос с trace ID
2. **UFO Service**: 
   - Server span: обработка `AnalyzeSighting`
   - Client span: `ufo.call_analysis` → Analysis Service
3. **Analysis Service**:
   - Server span: обработка `AnalyzeSighting`
   - Client span: `analysis.get_sighting` → UFO Service (`Get`)
   - Client span: `analysis.classify` → Classification Service
4. **Classification Service**:
   - Server span: обработка `ClassifyObject`
   - Internal span: `classification.analyze`

Все спаны связаны одним **trace ID** и передаются через gRPC metadata с помощью OpenTelemetry propagators.

## Быстрый старт

### 1. Запуск инфраструктуры

```bash
task up-all
```

### 2. Генерация конфигурации (если нужно)

```bash
task env:generate-all
```

### 3. Запуск сервисов

В отдельных терминалах:

```bash
# Terminal 1: UFO Service (с переменными из .env файла)
cd ufo
set -a && source ../deploy/compose/ufo/.env && set +a
go run cmd/grpc_server/main.go

# Terminal 2: Analysis Service
cd analysis
go run cmd/main.go

# Terminal 3: Classification Service
cd classification
go run cmd/main.go
```

### 4. Тестирование трейсинга

```bash
# Полный цикл тестирования (создание + анализ 3 объектов)
task tracing:test:full-cycle

# Проверка сервисов в Jaeger
task tracing:check-services

# Открытие Jaeger UI
task tracing:open-jaeger
```

## Команды тестирования

### Основные команды

```bash
# Полный цикл тестирования трейсинга
task tracing:test:full-cycle

# Создание отдельных типов объектов
task tracing:test:create-triangle
task tracing:test:create-saucer

# Проверка доступности сервисов в Jaeger
task tracing:check-services

# Открытие Jaeger UI в браузере
task tracing:open-jaeger
```

### Ручное тестирование

```bash
# Создание треугольного объекта
bin/grpcurl -plaintext -d '{"info": {"observed_at": "2024-01-15T21:00:00Z", "location": "Санкт-Петербург", "description": "Большой треугольный объект", "color": "красный", "sound": false, "duration_seconds": 600}}' localhost:50051 ufo.v1.UFOService/Create

# Создание классической тарелки
bin/grpcurl -plaintext -d '{"info": {"observed_at": "2024-01-15T22:00:00Z", "location": "Москва", "description": "Классическая летающая тарелка в форме диска", "color": "белый", "sound": true, "duration_seconds": 180}}' localhost:50051 ufo.v1.UFOService/Create

# Анализ наблюдения (замените UUID на реальный)
bin/grpcurl -plaintext -d '{"uuid": "YOUR_UUID_HERE"}' localhost:50051 ufo.v1.UFOService/AnalyzeSighting
```

## Просмотр трейсов в Jaeger

1. Откройте Jaeger UI: http://localhost:16686
2. Выберите сервис "ufo-service" в списке Services
3. Нажмите "Find Traces"
4. Выберите трейс для детального просмотра

### Что вы увидите в трейсах

Каждый вызов `AnalyzeSighting` создает трейс с 4 спанами:

1. **ufo.call_analysis** (UFO Service) - вызов Analysis Service
2. **analysis.get_sighting** (Analysis Service) - получение данных о наблюдении
3. **analysis.classify** (Analysis Service) - вызов Classification Service  
4. **classification.analyze** (Classification Service) - внутренняя логика классификации

## Примеры результатов классификации

### Треугольный объект
```json
{
  "analysisResult": "Анализ наблюдения: Объект классифицирован как unknown с уверенностью 0.55. Базовая классификация. Красный цвет может указывать на двигательную систему",
  "classification": "unknown", 
  "confidenceScore": 0.55
}
```

### Классическая тарелка
```json
{
  "analysisResult": "Анализ наблюдения: Объект классифицирован как classic_saucer с уверенностью 0.92. Классическая форма летающей тарелки. Белое свечение - распространенное явление",
  "classification": "classic_saucer",
  "confidenceScore": 0.92
}
```

### Сферический объект
```json
{
  "analysisResult": "Анализ наблюдения: Объект классифицирован как unknown с уверенностью 0.50. Базовая классификация",
  "classification": "unknown",
  "confidenceScore": 0.50
}
```

## Особенности реализации трейсинга

### Автоматическая пропагация контекста
- gRPC интерцепторы автоматически передают trace ID между сервисами
- Каждый запрос получает уникальный trace ID
- Все спаны в цепочке связаны одним trace ID

### Структура спанов
- **Server spans** - обработка входящих запросов
- **Client spans** - исходящие вызовы к другим сервисам
- **Internal spans** - внутренняя логика сервисов

### Метаданные трейсов
- Имена сервисов: ufo-service, analysis-service, classification-service
- Операции: методы gRPC и внутренние функции
- Атрибуты: параметры запросов, результаты обработки
- Ошибки: автоматическое логирование исключений

### Конфигурация трейсинга

Переменные окружения для UFO сервиса:
```bash
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_SERVICE_NAME=ufo-service
OTEL_ENVIRONMENT=development
OTEL_SERVICE_VERSION=1.0.0
```

## Устранение неполадок

### UFO сервис не появляется в Jaeger

1. Проверьте переменные окружения:
```bash
cat deploy/compose/ufo/.env | grep OTEL
```

2. Убедитесь, что UFO сервис запущен с правильными переменными:
```bash
cd ufo
set -a && source ../deploy/compose/ufo/.env && set +a
go run cmd/grpc_server/main.go
```

3. Проверьте логи OpenTelemetry Collector:
```bash
docker logs otel-collector --tail 20
```

### Проверка доступности сервисов

```bash
# Проверка сервисов в Jaeger
task tracing:check-services

# Проверка портов
lsof -i :50051 -i :50052 -i :50053 -i :16686

# Проверка контейнеров
docker ps
```

## Остановка системы

```bash
task down-all
```

## Полезные команды

```bash
# Просмотр доступных методов
bin/grpcurl -plaintext localhost:50051 list

# Описание сервиса
bin/grpcurl -plaintext localhost:50051 describe ufo.v1.UFOService

# Просмотр логов контейнеров
docker logs jaeger --tail 10
docker logs otel-collector --tail 10

# Просмотр всех задач Taskfile
task --list

# Генерация proto файлов
task proto:gen

# Форматирование кода
task format
```

## Архитектурные особенности

### Clean Architecture
- Каждый сервис использует clean architecture с разделением на слои
- Dependency injection через интерфейсы
- Конфигурация через переменные окружения

### OpenTelemetry интеграция
- Автоматическая инициализация трейсинга в каждом сервисе
- gRPC интерцепторы для автоматической пропагации контекста
- Корректное завершение трейсинга при остановке сервиса

### Микросервисная архитектура
- Каждый сервис имеет свою область ответственности
- Асинхронная коммуникация через gRPC
- Централизованное логирование и трейсинг 