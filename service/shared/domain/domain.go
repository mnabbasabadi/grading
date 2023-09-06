package domain

import (
	"sort"

	"github.com/google/uuid"
)

type (
	// ScaleType is the grade point average.
	ScaleType string
	// Grade ...
	Grade struct {
		StudentID uuid.UUID `db:"student_id"`
		CourseID  uuid.UUID `db:"course_id"`
		Grade     int       `db:"grade"`
	}

	GradeWithGPA struct {
		*Grade
		GPA string `db:"gpa"`
	}

	// Scale ...
	Scale struct {
		Min int    `db:"min"`
		GPA string `db:"gpa"`
	}
	// Scales ...
	Scales []Scale
)

const (
	DefaultScaleType ScaleType = "default"
)

// GetGPA is a method on Scales that returns the GPA for a given grade.
// if the grade is not found, it returns an empty string.
// scales should be sorted by Min in descending order.
func (s Scales) GetGPA(grade int) string {
	// Use binary search for faster lookup
	idx := sort.Search(len(s), func(i int) bool {
		return s[i].Min <= grade
	})

	// If we found a match
	if idx < len(s) {
		return s[idx].GPA
	}
	// Default GPA if not found
	return "N/A"
}
