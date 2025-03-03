package application

import (
	"context"

	"github.com/DimaKropachev/calculate-web-server/logger"
	"github.com/DimaKropachev/calculate-web-server/server/config"
	"github.com/DimaKropachev/calculate-web-server/server/internal/entities"
	"github.com/DimaKropachev/calculate-web-server/server/internal/orchestrator"
	"github.com/DimaKropachev/calculate-web-server/server/internal/service"
	"github.com/DimaKropachev/calculate-web-server/server/internal/storage/repository"
	"github.com/DimaKropachev/calculate-web-server/server/internal/transport/http/handlers"
	"github.com/DimaKropachev/calculate-web-server/server/internal/transport/http/router"
)

type Application struct {
	Config *config.Config
}

func NewServer(cfg *config.Config) *Application {
	return &Application{Config: cfg}
}

func (a *Application) Run() error {
	ctx := context.Background()
	ctx, err := logger.New(ctx)
	if err != nil {
		return err
	}

	idChan := make(chan string)

	storage := repository.NewExpressinsStorage()
	userService := service.NewService(storage)
	userHandler := handlers.NewUserHandler(userService, idChan)

	tq := entities.NewTasksQueue()
	rq := entities.NewResultsQueue()

	demonHandler := handlers.NewDemonHandler(tq, rq)

	orch := orchestrator.NewOrchestrator(storage, idChan, tq, rq)
	go orch.Start(ctx, a.Config.Timeouts)

	r := router.NewRouter(a.Config.Server, userHandler, demonHandler)
	r.Run(ctx)
	return nil
}
