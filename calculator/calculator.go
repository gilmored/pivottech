package calculator

import (
	"errors"
	"math"
)

func Add(a, b int) int {
	return a + b
}

func Subtract(a, b int) int {
	return a - b
}

func Multiply(a, b int) int {
	return a * b
}

func Divide(a, b int) (int, error) {
	err := errors.New("divide by zero")
	if b == 0 {
		return 0, err
	}

	return a / b, nil
}

func Power(x, y float64) float64 {
	return math.Pow(x, y)
}
