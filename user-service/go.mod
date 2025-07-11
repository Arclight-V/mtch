module user-service

go 1.23.8

require google.golang.org/grpc v1.72.0

require (
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)

require (
	github.com/go-passwd/validator v0.0.0-20250407044832-c284a2f4d990
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.33.0
	google.golang.org/protobuf v1.36.6
	proto v0.0.0
)

replace proto => ../proto
