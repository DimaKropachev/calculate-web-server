package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DimaKropachev/calculate-web-server/logger"
	"github.com/DimaKropachev/calculate-web-server/server/config"
	"github.com/DimaKropachev/calculate-web-server/server/internal/transport/http/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Router struct {
	config *config.ServerConfig
	Router *mux.Router

	UserHandleer handlers.UserHandler
	DemonHandler handlers.DemonHandler
}

func NewRouter(cfg *config.ServerConfig, uh *handlers.UserHandler, dh *handlers.DemonHandler) *Router {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./server/web/static/"))))

	r.HandleFunc("/internal/task", dh.GiveTask).Methods("GET")
	r.HandleFunc("/internal/task", dh.GetTask).Methods("POST")

	userApiRouter := r.PathPrefix("/api/v1").Subrouter()
	userApiRouter.Use(logger.LoggeringMiddleware)

	userApiRouter.HandleFunc("/find", handlers.FindExpression)
	userApiRouter.HandleFunc("/input", handlers.MainPage)
	userApiRouter.HandleFunc("/calculate", uh.AddExpr).Methods("POST")
	userApiRouter.HandleFunc("/expressions", uh.GetAllExprs).Methods("GET")
	userApiRouter.HandleFunc("/expressions/{id}", uh.GetExprById).Methods("GET")

	return &Router{config: cfg, Router: r}
}

func (r *Router) Run(ctx context.Context) {
	l := logger.GetLoggerFromCtx(ctx)
	l.Info(ctx, "server starting",
		zap.String("host", r.config.Host),
		zap.String("port", r.config.Port),
	)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", r.config.Host, r.config.Port), r.Router)
	if err != nil {
		l.Fatal(ctx, "server crashed",
			zap.String("host", r.config.Host),
			zap.String("port", r.config.Port),
			zap.String("error", err.Error()),
		)
	}
}
