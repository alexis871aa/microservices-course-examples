module github.com/olezhek28/microservices-course-examples/week_7/metrics/ufo

go 1.24.2

replace github.com/olezhek28/microservices-course-examples/week_7/metrics/platform => ../platform

replace github.com/olezhek28/microservices-course-examples/week_7/metrics/shared => ../shared

require (
	github.com/olezhek28/microservices-course-examples/week_7/metrics/platform v0.0.0-00010101000000-000000000000
	github.com/olezhek28/microservices-course-examples/week_7/metrics/shared v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.22.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.72.2
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
)
