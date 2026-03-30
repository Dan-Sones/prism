module data-cooking-service

go 1.25

require github.com/twmb/franz-go v1.20.7

require github.com/joho/godotenv v1.5.1

require (
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/grpc v1.79.2 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

require (
	github.com/ClickHouse/ch-go v0.71.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.43.0
	github.com/Dan-Sones/prismdbmodels v0.0.0 // indirect
	github.com/Dan-Sones/prismlogger v0.0.0
	github.com/Dan-Sones/prismmicrobatcher v0.0.0
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/paulmach/orb v0.12.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.26 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.12.0 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.40.0 // indirect
)

replace github.com/Dan-Sones/prismlogger => ../../libs/prismlogger

replace github.com/Dan-Sones/prismmicrobatcher => ../../libs/prismmicrobatcher

replace github.com/Dan-Sones/prismdbmodels => ../../libs/prismdbmodels
