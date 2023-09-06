//go:build integration
// +build integration

// Package integration contains integration tests for the service.
package integration

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/mnabbasabadi/grading/service/foundation/db"

	"github.com/jmoiron/sqlx"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"golang.org/x/exp/slog"

	migration "github.com/mnabbasabadi/grading/service/internal/storage/migration/postgres"
)

const (
	user     = "gnomock"
	password = "gnomick"
	database = "mydb"
)

func connectLocal() (*sqlx.DB, func() error, error) {

	p := postgres.Preset(
		postgres.WithUser(user, password),
		postgres.WithDatabase(database),
		postgres.WithTimezone("Europe/Paris"),
	)

	container, err := gnomock.Start(p)
	if err != nil {
		return nil, nil, err
	}

	opts := []db.Option{
		db.WithUser(user),
		db.WithPassword(password),
		db.WithHost(container.Host),
		db.WithPort(strconv.Itoa(container.DefaultPort())),
		db.WithDatabase(database),
		db.WithSSlMode(db.Disable),
	}
	pg, err := db.ConnectToPostgres(opts...)
	if err != nil {
		return nil, nil, err
	}
	return pg, func() error {
		return gnomock.Stop(container)
	}, nil
}
func connectHost() (*sqlx.DB, func() error, error) {
	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return nil, nil, fmt.Errorf("POSTGRES_USER not set")
	}
	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return nil, nil, fmt.Errorf("POSTGRES_PASSWORD not set")
	}
	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		return nil, nil, fmt.Errorf("POSTGRES_HOST not set")
	}
	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		return nil, nil, fmt.Errorf("POSTGRES_PORT not set")
	}
	database, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		return nil, nil, fmt.Errorf("POSTGRES_DB not set")
	}
	opts := []db.Option{
		db.WithUser(user),
		db.WithPassword(password),
		db.WithHost(host),
		db.WithPort(port),
		db.WithDatabase(database),
		db.WithSSlMode(db.Disable),
	}
	pg, err := db.ConnectToPostgres(opts...)
	if err != nil {
		return nil, nil, err
	}
	return pg, func() error {
		return pg.Close()
	}, nil
}
func getDB(logger *slog.Logger) (*sqlx.DB, func(), error) {
	isCli := isLocal()
	var (
		pg     *sqlx.DB
		err    error
		closer func() error
	)
	if isCli {
		pg, closer, err = connectLocal()
	} else {
		pg, closer, err = connectHost()
	}
	if err != nil {
		return nil, nil, err
	}

	if err := backoff.Retry(func() error {
		return pg.Ping()
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), 100)); err != nil {
		return nil, nil, err
	}

	if err := migration.GooseUP(pg.DB); err != nil {
		return nil, nil, err
	}
	return pg, func() {
		if err := closer(); err != nil {
			logger.With(err).Error("error closing database connection")
		}
	}, nil
}

func isLocal() bool {
	isCli, ok := os.LookupEnv("IS_CLI")
	if !ok {
		return true
	}
	return isCli == "true"
}
