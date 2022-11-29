package calculator_test

import (
	"github.com/gilmored/pivottech/calculator"
	"testing"
)

func TestCalculator(t *testing.T) {
	tests := map[string]struct {
		a, b, expect int
		x, y, want   float64
		err          string
		f            func(int, int) int
		fEr          func(int, int) (int, error)
		testPow      func(float64, float64) float64
	}{
		"Add":          {a: 2, b: 1, expect: 3, f: calculator.Add},
		"Subtract":     {a: 2, b: 1, expect: 1, f: calculator.Subtract},
		"Multiply":     {a: 2, b: 2, expect: 4, f: calculator.Multiply},
		"Divide":       {a: 4, b: 2, expect: 2, fEr: calculator.Divide},
		"DivideByZero": {a: 2, b: 0, expect: 0, fEr: calculator.Divide},
		"Power":        {x: 2, y: 3, want: 8, testPow: calculator.Power},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.f != nil {
				result := test.f(test.a, test.b)
				if result != test.expect {
					t.Errorf("got %d, expected %d", result, test.expect)
				}
				return
			}
			if test.fEr != nil {
				result, _ := test.fEr(test.a, test.b)
				if result != test.expect {
					t.Errorf("got %d, expected %d", result, test.expect)
				}
				return
			}
			if test.testPow != nil {
				result := test.testPow(test.x, test.y)
				if result != test.want {
					t.Errorf("got %v, expected %v", result, test.want)
				}
			}
		})
	}
}
