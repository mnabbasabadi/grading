package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mnabbasabadi/grading/service/shared/domain"
)

type (
	// Reader ...
	Reader struct {
		db *sqlx.DB
	}
)

// NewReader ...
func NewReader(db *sqlx.DB) Reader {
	return Reader{
		db: db,
	}
}

// language=postgresql
const getgpas = `select student_id, course_id, grade from grade order by created_at limit :limit offset :offset`

// language=postgresql
const totalgpas = `select count(*) from grade`

// GetGrades ...
func (r Reader) GetGrades(ctx context.Context, limit, offset int) ([]domain.Grade, int, error) {
	type params struct {
		Limit  int `db:"limit"`
		Offset int `db:"offset"`
	}
	p := params{
		Limit:  limit,
		Offset: offset,
	}

	stmt, err := r.db.PrepareNamedContext(ctx, getgpas)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	var gpas []domain.Grade
	if err := stmt.SelectContext(ctx, &gpas, p); err != nil {
		return nil, 0, fmt.Errorf("failed to get gpas: %w", err)
	}

	var total int
	if err := r.db.QueryRowxContext(ctx, totalgpas).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to get total: %w", err)
	}

	return gpas, total, nil
}

// language=postgresql
const getScale = `select min, gpa from scale where type=$1 order by min desc`

// GetScales ...
func (r Reader) GetScales(ctx context.Context, gpa domain.ScaleType) (domain.Scales, error) {
	var scales []domain.Scale
	err := r.db.SelectContext(ctx, &scales, getScale, gpa)
	if err != nil {
		return nil, err
	}
	if len(scales) == 0 {
		return nil, domain.ErrScaleNotFound
	}
	return scales, nil

}
