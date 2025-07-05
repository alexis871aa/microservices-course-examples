module github.com/alexis871aa/microservices-course-examples/week_1/workspace/grpc

go 1.24

require (
	github.com/alexis871aa/microservices-course-examples/week_1/workspace/shared v0.0.0-20250704154136-2ec32e31a857
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.35.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
)

replace github.com/alexis871aa/microservices-course-examples/week_1/workspace/shared => ../shared
