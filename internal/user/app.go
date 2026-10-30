package user

import (
	"github.com/gofiber/fiber/v2"
)

func NewApp(cfg Config, dbURL string) (*fiber.App, error) {
	db, err := openDatabase(dbURL)
	if err != nil {
		return nil, err
	}
	key, err := loadPrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		db.Close()
		return nil, err
	}
	_ = key
	app := fiber.New(fiber.Config{AppName: "ShopCore User Service", BodyLimit: 1 << 20})
	handler := &Handler{DB: db, PrivateKeyPath: cfg.PrivateKeyPath, AdminBootstrapKey: cfg.AdminBootstrapKey}
	app.Get("/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok", "service": "user"}) })
	app.Post("/api/v1/auth/bootstrap-admin", handler.BootstrapAdmin)
	app.Post("/api/v1/auth/register", handler.Register)
	app.Post("/api/v1/auth/login", handler.Login)
	app.Get("/api/v1/users/me", verifyUserToken(cfg.PublicKeyPath), handler.Me)
	return app, nil
}
