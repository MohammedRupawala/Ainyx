package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "Zero DOB",
			dob:      time.Time{},
			expected: 0,
		},
		{
			name:     "Future DOB",
			dob:      time.Now().Add(24 * time.Hour),
			expected: 0,
		},
		{
			name:     "Exact Birthday Today",
			dob:      time.Now().AddDate(-10, 0, 0),
			expected: 10,
		},
		{
			name:     "Birthday has not passed this year",
			dob:      time.Now().AddDate(-10, 0, 1),
			expected: 9,
		},
		{
			name:     "Birthday passed earlier this year",
			dob:      time.Now().AddDate(-10, 0, -1),
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateAge(tt.dob)
			if result != tt.expected {
				t.Errorf("calculateAge(%v) = %v; want %v", tt.dob, result, tt.expected)
			}
		})
	}
}
