module github.com/HaChains/infra

go 1.24.4

require (
	github.com/redis/go-redis/v9 v9.12.1
	github.com/rs/zerolog v1.34.0
	github.com/twmb/franz-go v1.20.4
	golang.org/x/crypto v0.43.0
	google.golang.org/protobuf v1.36.10
)

replace github.com/moodbased/go-lib => ../../moodbased/go-lib

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.12.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
)
