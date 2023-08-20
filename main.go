package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}

func main() {
	fmt.Printf("%f", Sqrt(2.0))
	//fmt.Println("Hello, world")
}
