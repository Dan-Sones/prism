module assignment-service

go 1.25

require (
	github.com/Dan-Sones/prismlogger v0.0.0
	github.com/go-chi/chi/v5 v5.2.4
	github.com/go-chi/cors v1.2.2
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.17.3
	github.com/stretchr/testify v1.11.1
	github.com/twmb/franz-go v1.20.6
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.12.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Dan-Sones/prismdbmodels => ../prismdbmodels

replace github.com/Dan-Sones/prismlogger => ../prismlogger
