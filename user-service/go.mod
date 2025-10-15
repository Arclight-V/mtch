module user-service

go 1.24.0

require google.golang.org/grpc v1.76.0

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.21.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff // indirect
)

require (
	github.com/Arclight-V/mtch v0.0.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/oklog/run v1.2.0
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.43.0
	google.golang.org/protobuf v1.36.10
)

replace github.com/Arclight-V/mtch => ..
