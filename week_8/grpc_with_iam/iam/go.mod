module github.com/olezhek28/microservices-course-examples/week_8/grpc_with_iam/iam

go 1.24.2

require (
	github.com/brianvoe/gofakeit/v7 v7.3.0
	github.com/envoyproxy/go-control-plane/envoy v1.32.4
	github.com/google/uuid v1.6.0
	github.com/olezhek28/microservices-course-examples/week_8/grpc_with_iam/shared v0.0.0-00010101000000-000000000000
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6
)

replace github.com/olezhek28/microservices-course-examples/week_8/grpc_with_iam/shared => ../shared

require (
	github.com/cncf/xds/go v0.0.0-20250501225837-2ac532fd4443 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250528174236-200df99c418a // indirect
)
