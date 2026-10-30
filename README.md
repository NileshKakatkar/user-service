# User Service

Fiber service for account registration, login, JWT issuing, and the authenticated user profile. PostgreSQL stores user accounts and passwords are hashed with bcrypt.

## Run

```sh
go mod download
go run ./cmd/user
```

Default port: `8081`. The service requires its own PostgreSQL database and the JWT key paths configured through environment variables.
