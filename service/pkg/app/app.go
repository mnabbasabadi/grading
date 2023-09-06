// Package app is the entry point of the service binary.
package app

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	kitHTTP "github.com/mnabbasabadi/grading/service/foundation/http"
	gradingAPI "github.com/mnabbasabadi/grading/service/internal/api/http"
	"github.com/mnabbasabadi/grading/service/internal/storage/rdbms/postgres"
	"github.com/mnabbasabadi/grading/service/internal/usecase"
	"golang.org/x/exp/slog"
)

// Environment ...
type Environment struct {
	logger *slog.Logger

	dbConn *sqlx.DB

	Logic usecase.Logic

	HTTPRegister kitHTTP.Registrar
}

// Params ...
type Params struct {
	//metrics *Metrics
	Logger       *slog.Logger
	HTTPRegister kitHTTP.Registrar

	// storage
	DB *sqlx.DB
}

// NewEnvironment ...
func NewEnvironment(ctx context.Context, params Params) *Environment {
	e := &Environment{
		//metrics: params.Metrics,
		logger: params.Logger,
		dbConn: params.DB,

		HTTPRegister: params.HTTPRegister,
	}
	e.Setup(ctx)
	return e
}

// Setup ...
func (e *Environment) Setup(_ context.Context) {
	repo := postgres.New(e.dbConn)
	logic := usecase.New(e.logger, repo)

	gradingHandler := gradingAPI.NewHandler(logic, e.logger)

	e.HTTPRegister(func(mux *http.ServeMux) {
		mw := []kitHTTP.Middleware{}
		// add metrics middleware
		//mw = append(mw, kitHTTP.Metrics(e.metrics))
		// add tracing middleware
		//mw = append(mw, kitHTTP.Tracing(e.tracer))
		// add logging middleware
		//mw = append(mw, kitHTTP.Logging(e.logger))
		// add authentication middleware
		//mw = append(mw, kitHTTP.Authentication(e.auth))
		h := kitHTTP.Chain(gradingHandler, mw...)
		mux.Handle("/", h)
	})
}

// Shutdown ...
func (e *Environment) Shutdown() {
	//e.metrics.Close()
}

// GetPostgresDB ...
func (e *Environment) GetPostgresDB() postgres.Repository {
	//return e.dbConn
	return nil
}
