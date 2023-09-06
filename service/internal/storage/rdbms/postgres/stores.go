// Package postgres is the implementation of the storage layer using PostgreSQL.
package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mnabbasabadi/grading/service/shared/domain"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=stores.go -package=postgres -destination=stores_mock.go Repository

type (
	// Repository is the interface that provides storage operations.
	Repository interface {
		GetGrades(context.Context, int, int) ([]domain.Grade, int, error)
		GetScales(context.Context, domain.ScaleType) (domain.Scales, error)
	}

	stores struct {
		Reader
	}
)

// New ...
func New(db *sqlx.DB) Repository {
	return &stores{
		Reader: NewReader(db),
	}
}
