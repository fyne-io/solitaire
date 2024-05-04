//go:generate fyne bundle --package=main -o data.go Icon.png

// Package main launches the solitaire app
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// show creates a new game and loads a table rendered in a new window.
func show(app fyne.App) {
	game := NewGame()
	table := NewTable(game)

	w := app.NewWindow("Solitaire")
	bar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			checkRestart(table, w)
		}))
	w.SetContent(container.NewBorder(bar, nil, nil, nil, table))
	w.Resize(fyne.NewSize(minWidth, minHeight))

	w.Show()
}

func checkRestart(t *Table, w fyne.Window) {
	dialog.ShowConfirm("New Game", "Start a new game?", func(ok bool) {
		if !ok {
			return
		}

		t.Restart()
	}, w)
}

func main() {
	a := app.New()
	a.SetIcon(resourceIconPng)
	a.Settings().SetTheme(newGameTheme())

	show(a)
	a.Run()
}
