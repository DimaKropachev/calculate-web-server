package orchestrator

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/DimaKropachev/calculate-web-server/logger"
	"github.com/DimaKropachev/calculate-web-server/server/config"
	"github.com/DimaKropachev/calculate-web-server/server/internal/entities"
	"github.com/DimaKropachev/calculate-web-server/server/internal/models"
	"github.com/DimaKropachev/calculate-web-server/server/pkg/calculate"
	"go.uber.org/zap"
)

type Service interface {
	Get(Id string) (*models.Expression, error)
}

type Orchestrator struct {
	service Service
	idChan  chan string
	tq      *entities.TasksQueue
	rq      *entities.ResultsQueue
}

func NewOrchestrator(service Service, idChan chan string, tq *entities.TasksQueue, rq *entities.ResultsQueue) *Orchestrator {
	return &Orchestrator{
		service: service,
		idChan:  idChan,
		tq:      tq,
		rq:      rq,
	}
}

func (o *Orchestrator) Start(ctx context.Context, timeCfg *config.OperationTime) {
	logger.GetLoggerFromCtx(ctx).Info(ctx, "orchestrator starting...")
	taskChan := make(chan *models.Task)

	go o.GetId(o.idChan, taskChan, ctx)
	HandleTask(taskChan, o.rq, o.tq, timeCfg)
}

func (o *Orchestrator) GetId(IDs <-chan string, tChan chan *models.Task, ctx context.Context) {
	for id := range IDs {
		go func() {
			expr, err := o.service.Get(id)
			if err != nil {
				logger.GetLoggerFromCtx(ctx).Info(ctx,
					"error when getting an expression by id",
					zap.String("id", id),
					zap.String("error", err.Error()),
				)
				return
			}

			expr.Status = "verify"
			err = calculate.CheckExpression(expr.Value)
			if err != nil {
				expr.Status = "error"
				expr.Error = err.Error()
				return
			}

			expr.Status = "calculate"
			tasks := Split(expr.Value, expr.ID)
			for _, task := range tasks {
				tChan <- task
			}
			if len(tasks) == 0 {
				expr.Status = "success"
				expr.Result = expr.Value
			} else {
				finalTaskId := tasks[len(tasks)-1].Id
				for {
					res, ok := o.rq.Get(finalTaskId)
					if ok {
						expr.Status = "success"
						expr.Result = fmt.Sprintf("%.2f", res.Result)
						break
					}
				}
			}
		}()
	}
}

func HandleTask(tasksChan chan *models.Task, rq *entities.ResultsQueue, tq *entities.TasksQueue, timeCfg *config.OperationTime) {
	mu := &sync.Mutex{}
	tasks := []*models.Task{}

	go func() {
		for task := range tasksChan {
			mu.Lock()
			tasks = append(tasks, task)
			mu.Unlock()
		}
	}()

	for {
		if len(tasks) == 0 {
			continue
		}
		mu.Lock()
		task := tasks[0]
		tasks = tasks[1:]
		mu.Unlock()

		if task.Arg1Task != "" {
			if res, ok := rq.Get(task.Arg1Task); ok {
				task.Arg1 = fmt.Sprintf("%.2f", res.Result)
				task.Arg1Task = ""
			}
		}
		if task.Arg2Task != "" {
			if res, ok := rq.Get(task.Arg2Task); ok {
				task.Arg2 = fmt.Sprintf("%.2f", res.Result)
				task.Arg2Task = ""
			}
		}

		if task.Arg1Task != "" || task.Arg2Task != "" {
			mu.Lock()
			tasks = append(tasks, task)
			mu.Unlock()
			continue
		}

		arg1, _ := strconv.ParseFloat(task.Arg1, 64)
		arg2, _ := strconv.ParseFloat(task.Arg2, 64)

		var timeOper time.Duration
		switch task.Oper {
		case "+":
			timeOper = timeCfg.Addition
		case "-":
			timeOper = timeCfg.Subtraction
		case "*":
			timeOper = timeCfg.Multiplication
		case "/":
			timeOper = timeCfg.Division
		}

		resTask := &models.FinalTask{
			ID:            task.Id,
			Arg1:          arg1,
			Arg2:          arg2,
			Operation:     task.Oper,
			OperationTime: timeOper,
		}
		tq.Add(resTask)
	}
}
