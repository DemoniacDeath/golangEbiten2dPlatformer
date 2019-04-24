package core

import "fmt"

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Times(scalar float64) Vector { return Vector{X: v.X * scalar, Y: v.Y * scalar} }

func (v Vector) Div(scalar float64) Vector { return Vector{X: v.X / scalar, Y: v.Y / scalar} }

func (v Vector) Plus(v2 Vector) Vector { return Vector{X: v.X + v2.X, Y: v.Y + v2.Y} }

func (v Vector) Minus(v2 Vector) Vector { return Vector{X: v.X - v2.X, Y: v.Y - v2.Y} }

func (v Vector) String() string {
	return fmt.Sprintf("{X: %f, Y: %f}", v.X, v.Y)
}
