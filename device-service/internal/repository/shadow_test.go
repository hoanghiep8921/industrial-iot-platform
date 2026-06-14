package repository

import (
	"testing"
)

func TestComputeDelta(t *testing.T) {
	tests := []struct {
		name     string
		desired  map[string]any
		reported map[string]any
		expected map[string]any
	}{
		{
			name:     "all match - no delta",
			desired:  map[string]any{"temp": 200},
			reported: map[string]any{"temp": 200},
			expected: map[string]any{},
		},
		{
			name:     "one differs",
			desired:  map[string]any{"temp": 200},
			reported: map[string]any{"temp": 195},
			expected: map[string]any{"temp": 200},
		},
		{
			name:     "missing key in reported",
			desired:  map[string]any{"temp": 200, "speed": 1500},
			reported: map[string]any{"temp": 200},
			expected: map[string]any{"speed": 1500},
		},
		{
			name:     "all differ",
			desired:  map[string]any{"temp": 200, "speed": 1500},
			reported: map[string]any{"temp": 100, "speed": 800},
			expected: map[string]any{"temp": 200, "speed": 1500},
		},
		{
			name:     "empty desired",
			desired:  map[string]any{},
			reported: map[string]any{"temp": 200},
			expected: map[string]any{},
		},
		{
			name:     "empty reported",
			desired:  map[string]any{"temp": 200},
			reported: map[string]any{},
			expected: map[string]any{"temp": 200},
		},
		{
			name:     "both empty",
			desired:  map[string]any{},
			reported: map[string]any{},
			expected: map[string]any{},
		},
		{
			name:     "string values match",
			desired:  map[string]any{"status": "RUNNING"},
			reported: map[string]any{"status": "RUNNING"},
			expected: map[string]any{},
		},
		{
			name:     "string values differ",
			desired:  map[string]any{"status": "STOPPED"},
			reported: map[string]any{"status": "RUNNING"},
			expected: map[string]any{"status": "STOPPED"},
		},
		{
			name:     "nested float values",
			desired:  map[string]any{"target": 200.5},
			reported: map[string]any{"target": 200.5},
			expected: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeDelta(tt.desired, tt.reported)
			if len(got) != len(tt.expected) {
				t.Errorf("computeDelta() len = %d, want %d; got=%v want=%v",
					len(got), len(tt.expected), got, tt.expected)
				return
			}
			for k, v := range tt.expected {
				if got[k] != v {
					t.Errorf("computeDelta()[%s] = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}
