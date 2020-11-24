// Package main launches the solitaire app
package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
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

	show(app)
	app.Run()
}
