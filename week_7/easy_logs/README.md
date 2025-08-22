# Easy Logs — централизованное логирование на OpenTelemetry + Elasticsearch + Kibana

Учебный, но полностью рабочий пример централизованного сбора логов: Go → OTLP gRPC → OpenTelemetry Collector → Elasticsearch → Kibana.

## 📋 Содержание

- [Концепция](#концепция)
- [Архитектура](#архитектура)
- [Компоненты](#компоненты)
- [OpenTelemetry Collector](#opentelemetry-collector)
- [Интеграция с Go](#интеграция-с-go)
- [Быстрый старт](#быстрый-старт)
- [Kibana](#kibana)
- [Мониторинг](#мониторинг)
- [Команды](#команды)

---

## Концепция

- Единый сбор логов со всех сервисов
- Структурированный JSON + единая схема (ECS)
- Доставка через агент (OTel Collector), батчирование и ретраи
- Метаданные сервиса: `service.name` выставляет само приложение; окружение — через коллектор

---

## Архитектура

```mermaid
graph TB
    subgraph "Go Applications"
        A[Note Service<br/>zap + Tee]
    end

    subgraph "Telemetry Layer"
        D[OpenTelemetry Collector<br/>OTLP gRPC :4317]
    end

    subgraph "Storage & Visualization"  
        E[Elasticsearch :9200]
        F[Kibana :5601]
    end

    A -->|OTLP gRPC| D
    D -->|Elasticsearch Exporter (REST)| E
    E -->|Query| F
```

Порты: 4317 (OTLP gRPC), 8888 (метрики коллектора), 9200 (Elasticsearch), 5601 (Kibana).

---

## Компоненты

- Elasticsearch: single-node, `xpack.security.enabled=false`, индекс `easy-logs`
- Kibana: автосоздание Data View `easy-logs*` (контейнер `kibana-init`), timeField `@timestamp`
- OpenTelemetry Collector: принимает логи по OTLP gRPC и шлёт в Elasticsearch (ECS mapping)

См. `docker-compose.yaml` (версии: Elasticsearch/Kibana 9.0.4, OTel Collector 0.123.0).

---

## OpenTelemetry Collector

Receivers:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
```

Processors (строго как в `otel-collector-config.yaml`):

```yaml
processors:
  resource:
    attributes:
      - key: deployment.environment
        value: dev
        action: insert  # service.name задаёт приложение; не переопределяем здесь
  batch:
    send_batch_size: 1000
    timeout: 5s
```

Exporters:

```yaml
exporters:
  elasticsearch:
    endpoints: ["http://elasticsearch:9200"]
    tls:
      insecure: true
    logs_index: "easy-logs"
    mapping:
      mode: ecs
```

Service/telemetry/extensions:

```yaml
extensions:
  health_check:
    endpoint: 0.0.0.0:13133

service:
  telemetry:
    metrics:
      address: 0.0.0.0:8888
  extensions: [health_check]
  pipelines:
    logs:
      receivers: [otlp]
      processors: [resource, batch]
      exporters: [elasticsearch]
```

---

## Интеграция с Go

В проекте используется zap c `zapcore.Tee`: запись одновременно в stdout и в OTLP через кастомный core.

Инициализация логгера в приложении (`cmd/grpc_server/main.go`):

```go
if err := applogger.Init("info", true, true); err != nil { // JSON + OTLP
    panic(err)
}
defer func() {
    _ = applogger.Sync()
    _ = applogger.Close()
}()
```

Ключевые моменты реализации (`internal/logger`):

- `createOTLPLogger` настраивает OTLP gRPC экспортер (endpoint `localhost:4317`) и `LoggerProvider` с `BatchProcessor` и ресурсами:

```go
resource.WithAttributes(
    semconv.ServiceName("note-service"),
    attribute.String("deployment.environment", "dev"),
)
```

- `SimpleOTLPCore` конвертирует `zap.Entry` в `otel/log.Record` и отправляет через OTLP с коротким таймаутом:

```go
record.SetSeverity(mapZapToOtelSeverity(entry.Level))
record.SetBody(otelLog.StringValue(entry.Message))
record.AddAttributes(encodeFieldsToAttrs(fields)...)
otlpLogger.Emit(ctx, record)
```

Сервис (`note_v1`) логирует события, напр.:

```go
applogger.Info(ctx, "Get request", zap.Int64("id", req.GetId()))
applogger.Error(ctx, "invalid id", zap.String("error", "id is empty"))
```

---

## Быстрый старт

1) Инфраструктура

```bash
task up
```

2) Приложение (отдельный терминал)

```bash
task run
```

gRPC сервер слушает `:50051`.

3) Генерация логов

```bash
task test:get
grpcurl -plaintext -d '{"id": 0}' localhost:50051 note_v1.NoteV1/Get
```

4) Kibana: откройте `http://localhost:5601` → Discover (Data View создаётся автоматически).

---

## Kibana

Полезные KQL-запросы:

```kql
# По уровню
event.severity_text: "ERROR"
event.severity_text: "INFO"

# По сообщению
message: "Get request"
message: *error* OR message: *failed*

# По сервису/окружению
service.name: "note-service"
deployment.environment: "dev"

# Комбинация
service.name: "note-service" AND event.severity_text: "ERROR"
```

Рекомендуемые колонки: `@timestamp`, `event.severity_text`, `message`, `service.name`, а также ваши поля из zap (`id`, `error`, ...).

---

## Мониторинг

- Метрики коллектора: `http://localhost:8888/metrics`
- Health check коллектора: `http://localhost:13133/health/status`

Prometheus-метрики, на которые смотреть:

```promql
otelcol_receiver_accepted_log_records_total{receiver="otlp"}
otelcol_exporter_sent_log_records_total{exporter="elasticsearch"}
otelcol_exporter_send_failed_log_records_total{exporter="elasticsearch"}
otelcol_processor_batch_batch_send_size_bucket{processor="batch"}
```

Диагностика:

- Индекс в ES: `curl localhost:9200/easy-logs/_count`
- Логи контейнеров: `docker logs elasticsearch|kibana|otel-collector|kibana-init`

---

## Команды

```bash
task up         # поднять инфраструктуру (Elasticsearch, Kibana, OTel Collector, init)
task down       # остановить и удалить контейнеры и volumes
task run        # запустить gRPC сервер с логированием
task proto:gen  # генерация protobuf (buf)
task format     # форматирование (gofumpt + gci)
task lint       # линтинг (golangci-lint)
task test:get   # тестовый gRPC запрос
task test:all   # все тесты
```

Поиск последних логов напрямую в ES:

```bash
curl "http://localhost:9200/easy-logs/_search?size=5&sort=@timestamp:desc" | jq .
```

---

## Ссылки

- OpenTelemetry: https://opentelemetry.io/docs/
- Elasticsearch: https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html
- Kibana: https://www.elastic.co/guide/en/kibana/current/index.html
- ECS: https://www.elastic.co/guide/en/ecs/current/index.html
- OpenTelemetry Go SDK: https://pkg.go.dev/go.opentelemetry.io/otel
