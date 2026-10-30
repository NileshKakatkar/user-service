# User Service

Fiber service for account registration, login, JWT issuing, and the authenticated user profile. PostgreSQL stores user accounts and passwords are hashed with bcrypt. The service also exposes the shared User gRPC contract on port `50051` for internal lookups.

## Run

```sh
go mod download
go run ./cmd/user
```

Default HTTP port: `8081`. Keep `shared/grpc` beside this service because its Go module uses the shared contracts repository through a local replace directive.
