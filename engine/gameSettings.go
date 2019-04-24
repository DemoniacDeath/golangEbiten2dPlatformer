package engine

import "../core"

type Settings struct {
	Title      string
	WindowSize core.Size
}

func NewSettings(title string, windowSize core.Size) *Settings {
	return &Settings{Title: title, WindowSize: windowSize}
}
