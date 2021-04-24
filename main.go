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

	w.Show()
}

func main() {
	app := app.New()
	app.SetIcon(resourceIconPng)
	app.Settings().SetTheme(newGameTheme())

	show(app)
	app.Run()
}
