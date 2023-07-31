package main

import (
	"context"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task/common"
	"task/internal/api/server"
)

func main() {

	logger := common.NewLogger()

	di, err := common.NewDIContainer()
	if err != nil {
		logger.Error("cannot create dependency injection container", slog.String("error", err.Error()))
		return
	}

	apiServer := server.NewServer(di)

	handler, err := apiServer.GetHTTPHandler(logger)
	if err != nil {
		logger.Error("cannot create http handler", slog.String("error", err.Error()))
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Addr:              di.Config.Address,
		Handler:           handler,
		ReadHeaderTimeout: di.Config.ReadHeaderTimeout,
		WriteTimeout:      di.Config.Timeout,
		IdleTimeout:       di.Config.IdleTimeout,
	}

	logger.Info("Server is running", slog.String("address", di.Config.Address))

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server stopped", slog.String("error", err.Error()))
	}

	<-done
	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), di.Config.ContextTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server", err)
		return
	}

	di.Pool.Close()

	logger.Info("server stopped")
}
