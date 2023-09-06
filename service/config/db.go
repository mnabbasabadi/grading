package config

import (
	"github.com/jmoiron/sqlx"
	"github.com/mnabbasabadi/grading/service/foundation/db"
	"golang.org/x/exp/slog"
)

// GetConnectionPQDB ...
func (c Config) GetConnectionPQDB(logger *slog.Logger) (*sqlx.DB, func(), error) {
	opts := []db.Option{
		db.WithUser(c.DB.User),
		db.WithPassword(c.DB.Password),
		db.WithHost(c.DB.Host),
		db.WithPort(c.DB.Port),
		db.WithDatabase(c.DB.DBName),
		db.WithSSlMode(c.DB.SslMode),
	}
	postgres, err := db.ConnectToPostgres(opts...)
	if err != nil {
		return nil, nil, err
	}
	return postgres, func() {
		if err := postgres.Close(); err != nil {
			logger.Error("error closing database connection", "err", err)
		}
	}, nil
}
