package main

import (
	"github.com/DimaKropachev/calculate-web-server/server/config"
	"github.com/DimaKropachev/calculate-web-server/server/internal/application"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	app := application.NewServer(cfg)
	panic(app.Run())
}
