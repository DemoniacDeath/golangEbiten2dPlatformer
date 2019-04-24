package core

type Rect struct {
	Center Vector
	Size   Size
}

func NewRect(x float64, y float64, width float64, height float64) Rect {
	return Rect{Vector{x, y}, Size{width, height}}
}
