package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	gradingAPI "github.com/mnabbasabadi/grading/api/v1"
	"github.com/mnabbasabadi/grading/service/internal/usecase"
	"github.com/mnabbasabadi/grading/service/shared/domain"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

func TestServer_GetGPA(t *testing.T) {
	testCases := map[string]struct {
		scaleType          domain.ScaleType
		setMock            func(m *usecase.MockLogic)
		expectedStatusCode int
		expectedGPAs       int
	}{
		"success": {
			scaleType: domain.ScaleType("4.0"),
			setMock: func(m *usecase.MockLogic) {
				m.EXPECT().GetGrades(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.GradeWithGPA{
					{
						Grade: &domain.Grade{
							StudentID: uuid.New(),
							CourseID:  uuid.New(),
							Grade:     25,
						},
						GPA: "F",
					},
					{
						Grade: &domain.Grade{
							StudentID: uuid.New(),
							CourseID:  uuid.New(),
							Grade:     33,
						},
						GPA: "D",
					},
					{
						Grade: &domain.Grade{
							StudentID: uuid.New(),
							CourseID:  uuid.New(),
							Grade:     43,
						},
						GPA: "C",
					},
				}, 100, nil)
			},
			expectedGPAs:       3,
			expectedStatusCode: http.StatusOK,
		},
		"failed to return logic- wrong scale type": {
			scaleType: domain.ScaleType("wrong"),
			setMock: func(m *usecase.MockLogic) {
				m.EXPECT().GetGrades(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.GradeWithGPA{}, 0, domain.ErrScaleNotFound)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		"failed to return logic- internal error": {
			scaleType: domain.ScaleType("wrong"),
			setMock: func(m *usecase.MockLogic) {
				m.EXPECT().GetGrades(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.GradeWithGPA{}, 0, errors.New("error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := usecase.NewMockLogic(ctrl)
			tc.setMock(mock)
			s := server{
				usecase: mock,
				logger:  logger,
			}
			req := httptest.NewRequest(http.MethodGet, "/students/gpa", nil)
			w := httptest.NewRecorder()
			ects := gradingAPI.GetGPAParamsScaleTypeECTS
			limit := 10
			offset := 0
			s.GetGPA(w, req, gradingAPI.GetGPAParams{
				ScaleType: &ects,
				Limit:     &limit,
				Offset:    &offset,
			})
			require.Equal(t, tc.expectedStatusCode, w.Code)

			defer func() {
				_ = w.Result().Body.Close()
			}()

			all, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)
			if w.Code == http.StatusOK {
				var responseBody gradingAPI.GradeList
				err = json.Unmarshal(all, &responseBody)
				require.NoError(t, err)
				require.Equal(t, tc.expectedGPAs, len(responseBody.Grades))
				require.Equal(t, 10, responseBody.Pagination.Limit)
				require.Equal(t, 0, responseBody.Pagination.Offset)
				require.Equal(t, 100, responseBody.Pagination.Total)

			}
		})
	}
}
