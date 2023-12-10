package main

import (
	"fmt"
	"math"
)

// Point vec.
type Point struct{ X, Y float64 }

// Pt creates a new point.
func Pt(x, y float64) Point { return Point{x, y} }

// String representation.
func (p Point) String() string {
	return fmt.Sprintf("(%g,%g)", math.Round(p.X*100)/100, math.Round(p.Y*100)/100)
}

// Scale the point.
func (p Point) Scale(fac float64) Point {
	return Point{p.X * fac, p.Y * fac}
}

// Add two points.
func (p Point) Add(p2 Point) Point {
	return Point{p.X + p2.X, p.Y + p2.Y}
}

// Mul two points.
func (p Point) Mul(p2 Point) Point {
	return Point{p.X * p2.X, p.Y * p2.Y}
}
