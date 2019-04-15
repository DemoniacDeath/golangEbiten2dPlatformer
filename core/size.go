package core

type Size struct {
	Width  float64
	Height float64
}

func (s Size) Times(scalar float64) Size { return Size{ Width: s.Width * scalar, Height: s.Height * scalar}}

func (s Size) Div(scalar float64) Size { return Size{ Width: s.Width / scalar, Height: s.Height / scalar}}

