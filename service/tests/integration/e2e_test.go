//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	gradingAPI "github.com/mnabbasabadi/grading/api/v1"
	"github.com/mnabbasabadi/grading/service/tests/support/client"
	"github.com/mnabbasabadi/grading/service/tests/support/storage/sqlt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	client   *client.GradeAPITestClient
	pgClient *sqlt.TestDAO
}

func (s *E2ETestSuite) SetupSuite() {
}

func (s *E2ETestSuite) TearDownSuite() {
}

func (s *E2ETestSuite) TestE2E() {
	s.T().Run("fail", func(t *testing.T) {
		t.Run("invalid scale type", func(t *testing.T) {
			ctx := context.Background()
			scaleType := gradingAPI.ScaleType("invalid")
			resp, err := s.client.GetGPA(ctx, scaleType, 10, 0)
			require.Error(t, err)
			require.Nil(t, resp.Grades)
		})
	})
	s.T().Run("success", func(t *testing.T) {
		ctx := context.Background()
		grade1 := gradingAPI.Grade{
			StudentId: uuid.NewString(),
			CourseId:  uuid.NewString(),
			Grade:     "2",
			Gpa:       "C",
		}
		grade2 := gradingAPI.Grade{
			StudentId: uuid.NewString(),
			CourseId:  uuid.NewString(),
			Grade:     "3",
			Gpa:       "B",
		}
		err := s.pgClient.InsertGrade(ctx, grade1.StudentId, grade1.CourseId, 2)
		require.NoError(t, err)
		err = s.pgClient.InsertGrade(ctx, grade2.StudentId, grade2.CourseId, 3)
		require.NoError(t, err)

		scaleType := gradingAPI.ScaleType("default")
		rsp, err := s.client.GetGPA(ctx, scaleType, 10, 0)
		require.NoError(t, err)
		require.Len(t, rsp.Grades, 2)
		require.Equal(t, grade1, rsp.Grades[0])
		require.Equal(t, grade2, rsp.Grades[1])

		require.Equal(t, rsp.Pagination, &gradingAPI.Pagination{
			Total:  2,
			Limit:  10,
			Offset: 0,
		})

	})
}
