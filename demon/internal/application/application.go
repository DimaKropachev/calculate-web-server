package application

import (
	"context"
	"time"

	"github.com/DimaKropachev/calculate-web-server/demon/config"
	"github.com/DimaKropachev/calculate-web-server/demon/internal/agent"
	"github.com/DimaKropachev/calculate-web-server/demon/internal/api"
	"github.com/DimaKropachev/calculate-web-server/demon/internal/model"
	"github.com/DimaKropachev/calculate-web-server/demon/internal/service"
	"github.com/DimaKropachev/calculate-web-server/logger"
	"go.uber.org/zap"
)

type Application struct {
	AgentConfig *config.Config
}

func NewAgent(cfg *config.Config) *Application {
	return &Application{
		AgentConfig: cfg,
	}
}

func (app *Application) Run() {
	ctx := context.Background()
	ctx, err := logger.New(ctx)
	if err != nil {
		panic(err)
	}

	api := api.NewApi("http://localhost:8081/internal/task")
	service := service.NewService(api)

	taskQueue := model.NewTaskQueue()

	go agent.StartAgents(ctx, app.AgentConfig.NumWorkers, taskQueue)
	go GetTasks(ctx, service, taskQueue.Tasks)
	GiveTasks(ctx, service, taskQueue.Results)
}

func GetTasks(ctx context.Context, service *service.Service, tasks chan<- *model.Task) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			task, err := service.Get()
			if task != nil {
				logger.GetLoggerFromCtx(ctx).Info(ctx,
					"task received",
					zap.String("task_id", task.ID),
				)
			}
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Info(ctx,
					"error when receiving the issue",
					zap.String("error", err.Error()),
				)
				continue
			}
			if task == nil {
				continue
			}
			tasks <- task
		}
	}
}

func GiveTasks(ctx context.Context, service *service.Service, results <-chan *model.Response) {
	for res := range results {
		select {
		case <-ctx.Done():
			return
		default:
			err := service.Give(res)
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Info(ctx,
					"error when sending the result",
					zap.String("error", err.Error()),
				)
			}
		}
	}
}
