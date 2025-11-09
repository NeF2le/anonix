module github.com/NeF2le/anonix/gateway

go 1.24.5

replace github.com/NeF2le/anonix/common => ../common

require (
	github.com/NeF2le/anonix/common v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/labstack/echo/v4 v4.13.4
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.9
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
