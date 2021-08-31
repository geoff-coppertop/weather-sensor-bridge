package math

import (
	"testing"
)

func TestRound(t *testing.T) {
	var tests = []struct {
		input     float64
		precision uint
		output    float64
	}{
		{1.0, 0, 1.0},
		{1.4, 0, 1.0},
		{1.5, 0, 2.0},
		{1.0, 1, 1.0},
		{1.4, 1, 1.4},
		{1.5, 1, 1.5},
	}

	for _, test := range tests {
		if test.output != Round(test.input, test.precision) {
			t.Errorf("expected %v, got %v", test.output, test.input)
		}
	}
}
