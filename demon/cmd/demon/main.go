package main

import (
	"github.com/DimaKropachev/calculate-web-server/demon/config"
	"github.com/DimaKropachev/calculate-web-server/demon/internal/application"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	agent := application.NewAgent(cfg)
	agent.Run()
}
