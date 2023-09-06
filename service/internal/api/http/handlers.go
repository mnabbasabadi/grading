// Package http is the entry point of the http binary.
package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	gradingAPI "github.com/mnabbasabadi/grading/api/v1"
	"github.com/mnabbasabadi/grading/service/shared/domain"

	"github.com/mnabbasabadi/grading/service/internal/usecase"

	"golang.org/x/exp/slog"
)

const (
	defaultLimit = 10
)

var _ gradingAPI.ServerInterface = new(server)

type (
	server struct {
		usecase usecase.Logic
		logger  *slog.Logger
	}
)

func (s server) GetLiveness(w http.ResponseWriter, _ *http.Request) {
	if err := respond(w, "OK", http.StatusOK); err != nil {
		s.logger.With(err).Error("while responding")
		return
	}
}

func (s server) GetReadiness(w http.ResponseWriter, _ *http.Request) {
	if err := respond(w, "OK", http.StatusOK); err != nil {
		s.logger.With(err).Error("while responding")
		return
	}
}

// GetGPA handles HTTP requests to get grades and calculate GPAs.
func (s server) GetGPA(w http.ResponseWriter, r *http.Request, params gradingAPI.GetGPAParams) {
	scaleType, limit, offset := s.parseParams(params)

	grades, total, err := s.usecase.GetGrades(r.Context(), scaleType, limit, offset)
	if err != nil {
		s.handleGradesError(w, err)
		return
	}

	response := s.prepareGradeResponse(grades, total, limit, offset)

	s.respond(w, response, http.StatusOK)
}

func (s server) parseParams(params gradingAPI.GetGPAParams) (domain.ScaleType, int, int) {
	var (
		scaleType domain.ScaleType
		limit     int
		offset    int
	)
	if params.ScaleType != nil {
		scaleType = domain.ScaleType(*params.ScaleType)
	}
	if params.Limit != nil {
		limit = *params.Limit
	}
	if limit == 0 {
		limit = defaultLimit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}
	return scaleType, limit, offset
}

func (s server) handleGradesError(w http.ResponseWriter, err error) {
	s.logger.Error("while getting grades", "error", err)
	if err == domain.ErrScaleNotFound {
		s.respondError(w, err, http.StatusBadRequest)
	} else {
		s.respondError(w, errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
	}
}

func (s server) prepareGradeResponse(grades []domain.GradeWithGPA, total, limit, offset int) gradingAPI.GradeList {
	var response gradingAPI.GradeList
	for _, grade := range grades {
		response.Grades = append(response.Grades, gradingAPI.Grade{
			CourseId:  grade.CourseID.String(),
			StudentId: grade.StudentID.String(),
			Grade:     fmt.Sprintf("%d", grade.Grade.Grade),
			Gpa:       grade.GPA,
		})
	}
	response.Pagination = &gradingAPI.Pagination{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	return response
}

// NewHandler returns a new http.Handler that implements the ServerInterface
func NewHandler(logic usecase.Logic, logger *slog.Logger) http.Handler {
	s := server{
		usecase: logic,
		logger:  logger,
	}

	options := gradingAPI.ChiServerOptions{
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.With(err).Error("error")
			s.respondError(w, err, http.StatusBadRequest)
		},
	}

	return gradingAPI.HandlerWithOptions(s, options)

}
func (s server) respond(w http.ResponseWriter, data any, statusCode int) {
	err := respond(w, data, statusCode)
	if err != nil {
		s.logger.With(err).Error("while responding")
		return
	}
}
