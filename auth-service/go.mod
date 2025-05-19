module github.com/Arclight-V/mtch/auth-service

go 1.24

require (
	goji.io v2.0.2+incompatible
	google.golang.org/grpc v1.72.0
)

require (
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	golang.org/x/crypto v0.33.0
	proto v0.0.0
)

replace proto => ../proto
