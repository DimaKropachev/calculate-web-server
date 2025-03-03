package agent

import (
	"context"
	"sync"
	"time"

	"github.com/DimaKropachev/calculate-web-server/demon/internal/model"
	"github.com/DimaKropachev/calculate-web-server/demon/pkg/calculate"
	"github.com/DimaKropachev/calculate-web-server/logger"
	"go.uber.org/zap"
)

type workerIDKey struct{}

func Agent(ctx context.Context, tasks <-chan *model.Task, results chan<- *model.Response, wg *sync.WaitGroup) {
	agentID := ctx.Value(workerIDKey{}).(int)

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"agent starting",
		zap.Int("agent_id", agentID),
	)

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.GetLoggerFromCtx(ctx).Info(ctx,
				"agent stop",
				zap.Int("agent_id", agentID),
			)
			return
		case task, ok := <-tasks:
			if ok {
				select {
				case <-time.After(task.OperationTime):
					logger.GetLoggerFromCtx(ctx).Info(ctx,
						"timeout",
						zap.String("id", task.ID),
						zap.Float64("arg1", task.Arg1),
						zap.Float64("arg2", task.Arg2),
						zap.String("operation", task.Operation),
						zap.Duration("oper_time", task.OperationTime),
					)
					continue
				default:
					now := time.Now()
					res, err := calculate.Calc(task.Arg1, task.Arg2, task.Operation)
					if err != nil {
						logger.GetLoggerFromCtx(ctx).Info(ctx,
							"agent error",
							zap.String("error", err.Error()),
						)
					}
					timeCalc := time.Since(now)
					time.Sleep(task.OperationTime - timeCalc)
					results <- &model.Response{
						ID:     task.ID,
						Result: res,
					}
				}
			}
		}
	}
}

func StartAgents(ctx context.Context, numWorkers int, taskQueue *model.TaskQueue) {
	wg := &sync.WaitGroup{}

	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		workerCtx := context.WithValue(ctx, workerIDKey{}, i)
		go Agent(workerCtx, taskQueue.Tasks, taskQueue.Results, wg)
	}

	wg.Wait()
}
