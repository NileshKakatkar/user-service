module github.com/nilesh/user-service

go 1.23

replace github.com/nilesh/shopcore-contracts => ../contracts

require (
	github.com/nilesh/shopcore-contracts v0.0.0
	google.golang.org/grpc v1.69.2
	github.com/gofiber/fiber/v2 v2.52.6
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/jackc/pgx/v5 v5.7.2
	golang.org/x/crypto v0.31.0
)
