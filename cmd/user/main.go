package main

import (
	"log"

	"github.com/nilesh/user-service/internal/user"
)

func main() {
	cfg := user.LoadConfig()
	app, err := user.NewApp(cfg, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
