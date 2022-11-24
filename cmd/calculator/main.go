package main

import (
	"fmt"
	"github.com/gilmored/pivottech/calculator"
)

func main() {
	fmt.Println(calculator.Add(1, 1))
	fmt.Println(calculator.Subtract(2, 1))
	fmt.Println(calculator.Multiply(2, 2))
	fmt.Println(calculator.Divide(6, 3))
	fmt.Println(calculator.Divide(4, 0))
}
