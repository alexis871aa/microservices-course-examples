module github.com/olezhek28/microservices-course-examples/week_1/workspace/grpc

go 1.24

require (
	github.com/google/uuid v1.6.0
	github.com/olezhek28/microservices-course-examples/week_1/workspace/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.72.2
	google.golang.org/protobuf v1.36.6
)

replace github.com/olezhek28/microservices-course-examples/week_1/workspace/shared => ../shared

require (
	github.com/google/go-cmp v0.7.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.35.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)
