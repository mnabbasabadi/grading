package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetGPA(t *testing.T) {
	testCases := map[string]struct {
		scales      Scales
		grade       int
		expectedGPA string
	}{
		"Grade Found in Middle": {
			scales: Scales{
				{Min: 90, GPA: "A"},
				{Min: 80, GPA: "B"},
				{Min: 70, GPA: "C"},
			},
			grade:       85,
			expectedGPA: "B",
		},
		"Grade Not Found": {
			scales: Scales{
				{Min: 90, GPA: "A"},
				{Min: 80, GPA: "B"},
				{Min: 70, GPA: "C"},
			},
			grade:       65,
			expectedGPA: "N/A",
		},
		"Exact Match": {
			scales: Scales{
				{Min: 90, GPA: "A"},
				{Min: 80, GPA: "B"},
				{Min: 70, GPA: "C"},
			},
			grade:       90,
			expectedGPA: "A",
		},
		"Grade Found at Beginning": {
			scales: Scales{
				{Min: 90, GPA: "A"},
				{Min: 80, GPA: "B"},
				{Min: 70, GPA: "C"},
			},
			grade:       95,
			expectedGPA: "A",
		},
		"Lowest Grade": {
			scales: Scales{
				{Min: 90, GPA: "A"},
				{Min: 80, GPA: "B"},
				{Min: 70, GPA: "C"},
			},
			grade:       60,
			expectedGPA: "N/A",
		},
		"Highest Grade": {
			scales: Scales{
				{Min: 90, GPA: "A"},
				{Min: 80, GPA: "B"},
				{Min: 70, GPA: "C"},
			},
			grade:       100,
			expectedGPA: "A",
		},
		"Empty Scale": {
			scales:      Scales{},
			grade:       85,
			expectedGPA: "N/A",
		},
		"Single Entry in Scale": {
			scales: Scales{
				{Min: 80, GPA: "B"},
			},
			grade:       85,
			expectedGPA: "B",
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			actualGPA := tc.scales.GetGPA(tc.grade)
			require.Equal(t, tc.expectedGPA, actualGPA)

		})
	}
}
