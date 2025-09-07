package tviewwrapper

import (
	"github.com/gdamore/tcell/v2"
)

type Style struct {
	BlurBackground  tcell.Color
	FocusBackground tcell.Color
}

func NewStyle() *Style {
	return &Style{
		BlurBackground:  tcell.ColorBlack,
		FocusBackground: tcell.Color234,

	}
}
