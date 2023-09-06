// Package usecase is the business usecase of the service. It is the place where the business rules are implemented.
package usecase

import (
	"context"
	"fmt"

	"github.com/mnabbasabadi/grading/service/internal/storage/rdbms/postgres"
	"github.com/mnabbasabadi/grading/service/shared/domain"
	"golang.org/x/exp/slog"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=logic.go -package=usecase -destination=logic_mock.go Logic

var _ Logic = new(controller)

type (
	// Logic is the interface that provides business usecase operations.
	Logic interface {
		GetGrades(ctx context.Context, scaleType domain.ScaleType, limit, offset int) ([]domain.GradeWithGPA, int, error)
	}
	controller struct {
		pg     postgres.Repository
		logger *slog.Logger
	}
)

// New returns a new Logic.
func New(logger *slog.Logger, pg postgres.Repository) Logic {
	return &controller{
		logger: logger,
		pg:     pg,
	}
}

// GetGrades fetches the grades and associates them with a GPA according to the given scaleType.
// It returns a slice of GradeWithGPA, the total number of grades, and any error encountered.
func (c *controller) GetGrades(ctx context.Context, scaleType domain.ScaleType, limit int, offset int) ([]domain.GradeWithGPA, int, error) {
	grades, total, err := c.fetchGrades(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fetching grades failed: %w", err)
	}
	if len(grades) == 0 {
		return nil, 0, fmt.Errorf("no grades provided")
	}

	scales, err := c.fetchScales(ctx, scaleType)
	if err != nil {
		return nil, 0, fmt.Errorf("fetching scales failed: %w", err)
	}

	gradesWithGPA, err := c.calculateGradesWithGPA(grades, scales)
	if err != nil {
		return nil, 0, fmt.Errorf("calculating grades with GPA failed: %w", err)
	}

	return gradesWithGPA, total, nil
}

func (c *controller) fetchGrades(ctx context.Context, limit int, offset int) ([]domain.Grade, int, error) {
	grades, total, err := c.pg.GetGrades(ctx, limit, offset)
	if err != nil {
		c.logger.Error("fetchGrades: failed to get grades", "error", err)
		return nil, 0, err
	}
	return grades, total, nil
}

func (c *controller) fetchScales(ctx context.Context, scaleType domain.ScaleType) (domain.Scales, error) {
	if scaleType == "" {
		scaleType = domain.DefaultScaleType
	}
	scales, err := c.pg.GetScales(ctx, scaleType)
	if err != nil {
		c.logger.Error("fetchScales: failed to get scales", "error", err)
		return nil, err
	}
	c.logger.With("scale", scales).Debug("fetched scales")
	return scales, nil
}

func (c *controller) calculateGradesWithGPA(grades []domain.Grade, scales domain.Scales) ([]domain.GradeWithGPA, error) {

	gradesWithGPA := make([]domain.GradeWithGPA, len(grades))
	for i, grade := range grades {
		gpa := scales.GetGPA(grade.Grade)
		gradesWithGPA[i] = domain.GradeWithGPA{
			Grade: &grades[i], // Directly use the address of the original slice element
			GPA:   gpa,
		}
	}
	return gradesWithGPA, nil
}
