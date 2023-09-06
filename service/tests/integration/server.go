//go:build integration
// +build integration

package integration

import (
	"time"

	"golang.org/x/exp/slog"

	kitHTTP "github.com/mnabbasabadi/grading/service/foundation/http"
)

func setupHTTPServer(shutdownTimeout, readTimeout, writeTimeout time.Duration, logger slog.Logger) kitHTTP.Server {
	serverOptions := []kitHTTP.ServerOption{
		kitHTTP.ShutdownTimeout(shutdownTimeout),
		kitHTTP.ReadTimeout(readTimeout),
		kitHTTP.WriteTimeout(writeTimeout),
	}
	return kitHTTP.NewServer(logger, serverOptions...)
}
