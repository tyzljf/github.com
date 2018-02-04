package main

import (
	"math"
	"fmt"
)

type Abser interface {
	Abs() float64
}

/*
	Vertex implement
 */
type Vertex struct {
	X, Y float64
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

/*
	MyFloat implement
 */
type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}


func main() {
	var a Abser

	v := Vertex{3, 4}
	f := MyFloat(-math.Sqrt2)

	a = &v
	fmt.Println(a.Abs())

	a = &f
	fmt.Println(a.Abs())

}

