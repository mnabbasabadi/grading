// Package integration is the integration test suite of the service.

//go:build integration
// +build integration

package integration

import (
	"context"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	gradingAPI "github.com/mnabbasabadi/grading/api/v1"
	"github.com/mnabbasabadi/grading/service/pkg/app"
	"github.com/mnabbasabadi/grading/service/tests/support/client"
	"github.com/mnabbasabadi/grading/service/tests/support/storage/sqlt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slog"
)

const (
	host = "localhost"
	port = "8080"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	ctx := context.TODO()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dbConn, dbCloser, err := getDB(logger)
	require.NoError(t, err)
	defer dbCloser()

	httpServer := setupHTTPServer(5*time.Second, 5*time.Second, 5*time.Second, *logger)
	defer httpServer.Stop()

	params := app.Params{
		Logger:       logger,
		DB:           dbConn,
		HTTPRegister: httpServer.Register,
	}

	env := app.NewEnvironment(ctx, params)
	defer env.Shutdown()

	addr := net.JoinHostPort(host, port)
	httpServer.Start(addr, nil)

	testClient, err := client.NewGradingAPITestClient(addr, gradingAPI.WithHTTPClient(http.DefaultClient))
	require.NoError(t, err)

	suites := map[string]suite.TestingSuite{
		"E2E": &E2ETestSuite{
			client:   testClient,
			pgClient: sqlt.NewTestDAO(dbConn),
		},
	}
	for _, s := range suites {
		suite.Run(t, s)
	}
}
