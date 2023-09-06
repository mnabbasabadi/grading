package usecase

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/mnabbasabadi/grading/service/internal/storage/rdbms/postgres"
	"github.com/mnabbasabadi/grading/service/shared/domain"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

func TestController_GetGrades(t *testing.T) {

	testCases := map[string]struct {
		gpa         domain.ScaleType
		setMock     func(m *postgres.MockRepository)
		expectedGPA []string
		wantErr     bool
	}{
		"success": {
			gpa: domain.ScaleType("4.0"),
			setMock: func(m *postgres.MockRepository) {
				m.EXPECT().GetGrades(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Grade{
					{
						StudentID: uuid.New(),
						CourseID:  uuid.New(),
						Grade:     25,
					},
					{
						StudentID: uuid.New(),
						CourseID:  uuid.New(),
						Grade:     33,
					},
					{
						StudentID: uuid.New(),
						CourseID:  uuid.New(),
						Grade:     43,
					},
				}, 3, nil)
				m.EXPECT().GetScales(gomock.Any(), gomock.Any()).Return(domain.Scales{
					{
						Min: 40,
						GPA: "C",
					},
					{
						Min: 30,
						GPA: "D",
					},
					{
						Min: 20,
						GPA: "F",
					},
				}, nil)
			},
			expectedGPA: []string{"F", "D", "C"},
		},
		"success with default gpa": {
			gpa: domain.ScaleType(""),
			setMock: func(m *postgres.MockRepository) {
				m.EXPECT().GetGrades(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Grade{
					{
						StudentID: uuid.New(),
						CourseID:  uuid.New(),
						Grade:     25,
					},
					{
						StudentID: uuid.New(),
						CourseID:  uuid.New(),
						Grade:     33,
					},
					{
						StudentID: uuid.New(),
						CourseID:  uuid.New(),
						Grade:     43,
					},
				}, 3, nil)
				m.EXPECT().GetScales(gomock.Any(), gomock.Any()).Return(domain.Scales{
					{
						Min: 40,
						GPA: "C",
					},
					{
						Min: 30,
						GPA: "D",
					},
					{
						Min: 20,
						GPA: "F",
					},
				}, nil)
			},
			expectedGPA: []string{"F", "D", "C"},
		},
		"fail to get grades": {
			gpa: domain.ScaleType(""),
			setMock: func(m *postgres.MockRepository) {
				m.EXPECT().GetGrades(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, 0, errors.New("error"))
			},
			wantErr: true,
		},
		"fail to get scales": {
			gpa: domain.ScaleType(""),
			setMock: func(m *postgres.MockRepository) {
				m.EXPECT().GetGrades(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Grade{
					{
						StudentID: uuid.New(),
						CourseID:  uuid.New(),
						Grade:     25,
					},
					{
						StudentID: uuid.New(),
						CourseID:  uuid.New(),
						Grade:     33,
					},
					{
						StudentID: uuid.New(),
						CourseID:  uuid.New(),
						Grade:     43,
					},
				}, 3, nil)
				m.EXPECT().GetScales(gomock.Any(), gomock.Any()).Return(nil, domain.ErrScaleNotFound)
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := postgres.NewMockRepository(ctrl)
			if tc.setMock != nil {
				tc.setMock(m)
			}
			c := controller{
				pg:     m,
				logger: logger,
			}
			grades, total, err := c.GetGrades(context.TODO(), tc.gpa, 10, 0)
			require.Equal(t, tc.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, 3, total)
			for i, grade := range grades {
				require.Equal(t, tc.expectedGPA[i], grade.GPA)
			}
		})
	}
}
