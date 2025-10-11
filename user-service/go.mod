module user-service

go 1.24.0

require google.golang.org/grpc v1.76.0

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.21.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)

require (
	config v0.0.0
	github.com/Arclight-V/mtch/pkg/prober v0.0.0-00010101000000-000000000000
	github.com/Arclight-V/mtch/pkg/signaler v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/oklog/run v1.2.0
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.40.0
	google.golang.org/protobuf v1.36.6
	proto v0.0.0
)

replace proto => ../proto

replace config => ./../pkg/platform/config

replace github.com/Arclight-V/mtch/pkg/signaler => ./../pkg/signaler

replace github.com/Arclight-V/mtch/pkg/prober => ./../pkg/prober
