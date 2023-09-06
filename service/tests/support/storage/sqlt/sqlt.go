package sqlt

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mnabbasabadi/grading/service/internal/storage/rdbms/postgres"
)

// TestDAO ...
type TestDAO struct {
	postgres.Repository
	db *sqlx.DB
}

// NewTestDAO ...
func NewTestDAO(db *sqlx.DB) *TestDAO {
	return &TestDAO{
		db:         db,
		Repository: postgres.New(db),
	}
}

// language=postgresql
const insertgrade = `INSERT INTO grade (student_id, course_id, grade) VALUES ($1, $2, $3)`

// InsertGrade ...
func (t *TestDAO) InsertGrade(ctx context.Context, student_id, course_id string, grade int) error {
	if _, err := t.db.ExecContext(ctx, insertgrade, student_id, course_id, grade); err != nil {
		return err
	}
	return nil
}
