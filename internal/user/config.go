package user

import "os"

type Config struct {
	Port              string
	GRPCPort          string
	DatabaseURL       string
	PrivateKeyPath    string
	PublicKeyPath     string
	AdminBootstrapKey string
}

func LoadConfig() Config {
	return Config{
		Port:              env("USER_SERVICE_PORT", "8081"),
		GRPCPort:          env("USER_GRPC_PORT", "50051"),
		DatabaseURL:       databaseURL("shopcore_users"),
		PrivateKeyPath:    env("JWT_PRIVATE_KEY_PATH", "./secrets/private.pem"),
		PublicKeyPath:     env("JWT_PUBLIC_KEY_PATH", "./secrets/public.pem"),
		AdminBootstrapKey: env("ADMIN_BOOTSTRAP_KEY", ""),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func databaseURL(name string) string {
	return "postgres://" + env("POSTGRES_USER", "shopcore") + ":" + env("POSTGRES_PASSWORD", "shopcore") +
		"@" + env("POSTGRES_HOST", "localhost") + ":" + env("POSTGRES_PORT", "5432") + "/" + name +
		"?sslmode=" + env("POSTGRES_SSLMODE", "disable")
}
