module github.com/olezhek28/microservices-course-examples/week_8/grpc_with_iam/ufo

go 1.24.2

replace github.com/olezhek28/microservices-course-examples/week_8/grpc_with_iam/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/olezhek28/microservices-course-examples/week_8/grpc_with_iam/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/samber/lo v1.51.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)
