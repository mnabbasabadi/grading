// Package main is the entry point of the grading binary.
package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mnabbasabadi/grading/service/config"
	"github.com/mnabbasabadi/grading/service/internal/storage/migration/postgres"
	"github.com/mnabbasabadi/grading/service/pkg/app"
	"golang.org/x/exp/slog"

	kitHTTP "github.com/mnabbasabadi/grading/service/foundation/http"
)

func main() {
	ctx := context.Background()
	// Load config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Start server
	serverErrors := make(chan error, 1)

	defer func() {
		if err != nil {
			logger.Error("run error", "err", err)
		}
	}()

	defer config.RecoverAndLogPanic(logger)

	httpServer := setupHTTPServer(5*time.Second, 5*time.Second, 5*time.Second, *logger)

	dbConn, dbCloser, err := cfg.GetConnectionPQDB(logger)
	if err != nil {
		logger.Error("error setting up database connection", "err", err)
		os.Exit(1)
	}
	defer dbCloser()

	// migrate database
	if err := postgres.GooseUP(dbConn.DB); err != nil {
		logger.Error("error running migrations", "err", err)
		os.Exit(1)
	}

	params := app.Params{
		Logger:       logger,
		DB:           dbConn,
		HTTPRegister: httpServer.Register,
	}

	env := app.NewEnvironment(ctx, params)

	logger.Info("setup complete, starting server")
	startHTTPServer(httpServer, cfg.Host, cfg.Port, serverErrors)
	shutdown := listenForShutdown()
	select {
	case err := <-serverErrors:
		logger.Error("error starting server", "err", err)
		os.Exit(1)
	case sig := <-shutdown:
		logger.Error("caught signal", "signal", sig)
		env.Shutdown()
		httpServer.Stop()
	}
}

// ListenForShutdown creates a channel and subscribes to specific signals to trigger a shutdown of the service.
func listenForShutdown() chan os.Signal {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	return shutdown
}

func setupHTTPServer(shutdownTimeout, readTimeout, writeTimeout time.Duration, logger slog.Logger) kitHTTP.Server {
	serverOptions := []kitHTTP.ServerOption{
		kitHTTP.ShutdownTimeout(shutdownTimeout),
		kitHTTP.ReadTimeout(readTimeout),
		kitHTTP.WriteTimeout(writeTimeout),
	}
	return kitHTTP.NewServer(logger, serverOptions...)
}

func startHTTPServer(httpServer kitHTTP.Server, host, port string, serverErrors chan<- error) {
	addr := net.JoinHostPort(host, port)
	httpServer.Start(addr, serverErrors)
}
