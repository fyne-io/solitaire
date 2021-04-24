// Package main launches the solitaire app
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// show creates a new game and loads a table rendered in a new window.
func show(app fyne.App) {
	game := NewGame()

	w := app.NewWindow("Solitaire")
	w.SetPadded(false)
	w.SetContent(NewTable(game))
	w.Resize(fyne.NewSize(minWidth, minHeight))

	w.Show()
}

func main() {
	a := app.New()
	a.SetIcon(resourceIconPng)
	a.Settings().SetTheme(newGameTheme())

	show(a)
	a.Run()
}
