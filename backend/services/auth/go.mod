module github.com/williamsebastianliman/WEB-WS-242/services/auth

replace github.com/williamsebastianliman/WEB-WS-242/services/auth => ../auth

replace github.com/williamsebastianliman/WEB-WS-242/proto/gen/go => ../../proto/gen/go

go 1.23.4

require (
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/google/uuid v1.6.0
	github.com/redis/go-redis/v9 v9.10.0
	github.com/streadway/amqp v1.1.0
	github.com/williamsebastianliman/WEB-WS-242/proto/gen/go v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.39.0
	google.golang.org/grpc v1.73.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
