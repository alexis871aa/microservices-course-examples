module github.com/olezhek28/microservices-course-examples/week_7/tracing/classification

go 1.24.2

require (
	github.com/olezhek28/microservices-course-examples/week_7/tracing/platform v0.0.0-00010101000000-000000000000
	github.com/olezhek28/microservices-course-examples/week_7/tracing/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.72.2
)

require (
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.36.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/olezhek28/microservices-course-examples/week_7/tracing/platform => ../platform

replace github.com/olezhek28/microservices-course-examples/week_7/tracing/shared => ../shared
