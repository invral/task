package main

import (
	"golang.org/x/exp/slog"
	"net/http"
	"task/common"
	"task/internal/api/server"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	logger := common.NewLogger()

	di, err := common.NewDIContainer()
	if err != nil {
		logger.Error("cannot create dependency injection container", slog.String("error", err.Error()))
		return
	}
	//_ = di

	apiServer := server.NewServer(di)

	handler, err := apiServer.GetHTTPHandler(logger)
	if err != nil {
		logger.Error("cannot create http handler", slog.String("error", err.Error()))
	}

	server := http.Server{
		Addr:              di.Config.Address,
		Handler:           handler,
		ReadHeaderTimeout: di.Config.ReadHeaderTimeout,
		WriteTimeout:      di.Config.Timeout,
		IdleTimeout:       di.Config.IdleTimeout,
	}

	logger.Info("Server is running", slog.String("address", di.Config.Address))

	if err := server.ListenAndServe(); err != nil {
		logger.Error("server stopped", slog.String("error", err.Error()))
	}
}
