package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type gameTheme struct {
	fyne.Theme
}

func newGameTheme() fyne.Theme {
	return &gameTheme{theme.DefaultTheme()}
}

func (g *gameTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.RGBA{R: 0x07, G: 0x63, B: 0x24, A: 0xff}
	case theme.ColorNameSeparator:
		return color.RGBA{R: 0x02, G: 0x52, B: 0x10, A: 0xff}
	}

	return g.Theme.Color(n, v)
}
