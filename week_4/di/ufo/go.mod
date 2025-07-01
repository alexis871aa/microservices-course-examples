module github.com/olezhek28/microservices-course-examples/week_4/di/ufo

go 1.24

replace github.com/olezhek28/microservices-course-examples/week_4/di/shared => ../shared

replace github.com/olezhek28/microservices-course-examples/week_4/di/platform => ../platform

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/olezhek28/microservices-course-examples/week_4/di/platform v0.0.0-00010101000000-000000000000
	github.com/olezhek28/microservices-course-examples/week_4/di/shared v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/samber/lo v1.50.0
	go.mongodb.org/mongo-driver v1.17.3
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.72.2
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/golang/snappy v1.0.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
)
