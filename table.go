package main

import (
	"math"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fyne-io/solitaire/faces"
)

// Table represents the rendering of a game in progress
type Table struct {
	widget.BaseWidget

	game     *Game
	selected *Card

	float       *canvas.Image
	floatSource *canvas.Image
	floatPos    fyne.Position
	shuffle     *widget.ToolbarAction

	findCard func(fyne.Position) (*Card, *canvas.Image)
	stackPos func(int) fyne.Position
}

// CreateRenderer gets the widget renderer for this table - internal use only
func (t *Table) CreateRenderer() fyne.WidgetRenderer {
	return newTableRender(t)
}

// find card from an image, easier than keeping them in sync
func (t *Table) cardForPos(pos *canvas.Image) *Card {
	deck := NewSortedDeck()

	for i, face := range deck.Cards {
		if face.Face() == pos.Resource {
			card := NewCard((i%13)+1, Suit(math.Floor(float64(i)/13)))
			card.FaceUp = true // we know this as we checked the face
			return card
		}
	}

	return nil
}

func cardEquals(card1, card2 *Card) bool {
	if card1 == nil || card2 == nil {
		return card1 == nil && card2 == nil
	}

	return card1.Value == card2.Value && card1.Suit == card2.Suit
}

func (t *Table) cardTapped(cardPos *canvas.Image, pos fyne.Position, move func()) bool {
	if !withinCardBounds(cardPos, pos) {
		return false
	}

	card := t.cardForPos(cardPos)
	if cardPos.Resource != faces.ForSpace() && (card == nil || !card.FaceUp) {
		t.selected = nil
		t.Refresh()

		return true
	}

	if t.selected == nil {
		t.selected = card
	} else {
		if cardEquals(t.selected, card) {
			t.game.AutoBuild(card)
		} else {
			if move != nil {
				move()
			}
		}

		t.selected = nil
	}

	t.Refresh()
	return true
}

func (t *Table) checkStackTapped(render *stackRender, stack *Stack, pos fyne.Position) bool {
	for i := len(stack.Cards) - 1; i >= 0; i-- {
		//		card := stack.Cards[i]

		if t.cardTapped(render.cards[i], pos, func() {
			t.game.MoveCardToStack(stack, t.selected)
		}) {
			return true
		}
	}

	return t.cardTapped(render.cards[0], pos, func() {
		t.game.MoveCardToStack(stack, t.selected)
	})
}

func (t *Table) Restart() {
	oldWin := t.game.OnWin
	t.game = NewGame()
	t.game.OnWin = oldWin
	t.shuffle.Enable()

	test.WidgetRenderer(t).(*tableRender).game = t.game
	t.Refresh()
}

// Dragged is called when the user drags on the table widget
func (t *Table) Dragged(event *fyne.DragEvent) {
	t.floatPos = event.Position
	if !t.float.Hidden { // existing drag
		t.float.Move(t.float.Position().Add(event.Dragged))
		return
	}

	if t.selected != nil {
		return
	}

	card, source := t.findCard(event.Position)
	if card == nil {
		return
	}
	if !card.FaceUp { // only drag visible cards
		return
	}

	t.floatSource = source
	t.float.Resource = source.Resource
	t.selected = card
	source.Resource = nil
	source.Image = nil
	source.Refresh()
	t.float.Refresh()
	t.float.Move(source.Position())
	t.float.Show()
}

// DragEnd is called when the user stops dragging on the table widget
func (t *Table) DragEnd() {
	t.float.Hide()
	if t.dropCard(t.floatPos) {
		return
	}

	if t.floatSource != nil {
		t.floatSource.Resource = t.float.Resource
		t.floatSource.Refresh()
	}
}

// Tapped is called when the user taps the table widget
func (t *Table) Tapped(event *fyne.PointEvent) {
	render := test.WidgetRenderer(t).(*tableRender)

	if withinCardBounds(render.deck, event.Position) {
		t.selected = nil
		t.game.DrawThree()

		if len(t.game.Drawn.Cards) == 0 {
			t.shuffle.Enable()
		} else {
			t.shuffle.Disable()
		}

		render.Refresh()
		return
	}

	t.dropCard(event.Position)
}

func (t *Table) dropCard(pos fyne.Position) bool {
	render := test.WidgetRenderer(t).(*tableRender)

	if t.game.Draw3 != nil {
		if t.cardTapped(render.pile3, pos, nil) {
			return true
		}
	} else if t.game.Draw2 != nil {
		if t.cardTapped(render.pile2, pos, nil) {
			return true
		}
	} else if t.game.Draw1 != nil {
		if t.cardTapped(render.pile1, pos, nil) {
			return true
		}
	}

	if t.cardTapped(render.build1, pos, func() {
		t.game.MoveCardToBuild(t.game.Build1, t.selected)
	}) {
		return true
	} else if t.cardTapped(render.build2, pos, func() {
		t.game.MoveCardToBuild(t.game.Build2, t.selected)
	}) {
		return true
	} else if t.cardTapped(render.build3, pos, func() {
		t.game.MoveCardToBuild(t.game.Build3, t.selected)
	}) {
		return true
	} else if t.cardTapped(render.build4, pos, func() {
		t.game.MoveCardToBuild(t.game.Build4, t.selected)
	}) {
		return true
	}

	if t.checkStackTapped(render.stack1, t.game.Stack1, pos) {
		return true
	} else if t.checkStackTapped(render.stack2, t.game.Stack2, pos) {
		return true
	} else if t.checkStackTapped(render.stack3, t.game.Stack3, pos) {
		return true
	} else if t.checkStackTapped(render.stack4, t.game.Stack4, pos) {
		return true
	} else if t.checkStackTapped(render.stack5, t.game.Stack5, pos) {
		return true
	} else if t.checkStackTapped(render.stack6, t.game.Stack6, pos) {
		return true
	} else if t.checkStackTapped(render.stack7, t.game.Stack7, pos) {
		return true
	}

	t.selected = nil // clicked elsewhere
	t.Refresh()
	return false
}

func (t *Table) finishAnimation() {
	anim := container.NewWithoutLayout()
	anim.Resize(t.Size())
	fyne.CurrentApp().Driver().CanvasForObject(t).Overlays().Add(anim)
	wg := &sync.WaitGroup{}
	for i := ValueKing; i > 0; i-- {
		for j, p := range []*Stack{t.game.Build1, t.game.Build2, t.game.Build3, t.game.Build4} {
			card := p.Pop()
			pos := t.stackPos(j)

			wg.Add(1)
			off := fyne.Delta{DX: -2, DY: 1}
			switch j {
			case 0:
				off.DX = 2
			case 1:
				off.DX = 2
				off.DY = -1
			case 2:
				off.DY = -1
			}
			image := t.startCardAnimation(card, pos, off, wg)

			anim.Objects = append([]fyne.CanvasObject{image}, anim.Objects...)
			t.Refresh()
			time.Sleep(time.Second / 6)
		}
	}

	wg.Wait()
	anim.Objects = []fyne.CanvasObject{}
	fyne.CurrentApp().Driver().CanvasForObject(t).Overlays().Remove(anim)
}

func (t *Table) startCardAnimation(card *Card, pos fyne.Position, off fyne.Delta, wg *sync.WaitGroup) fyne.CanvasObject {
	bounds := t.Size()
	pad := theme.Padding()
	i := canvas.NewImageFromResource(faces.ForCard(card.Value, int(card.Suit)))
	i.Resize(cardSize)
	i.Move(pos)

	go func() {
		for pos.X > -cardSize.Width-pad && pos.Y > -cardSize.Height-pad && pos.X < bounds.Width && pos.Y < bounds.Height {
			pos = pos.Add(off)
			i.Move(pos)

			time.Sleep(time.Millisecond * 12)
		}

		i.Hide()
		wg.Done()
	}()

	return i
}

// NewTable creates a new table widget for the specified game
func NewTable(g *Game) *Table {
	table := &Table{game: g}
	table.ExtendBaseWidget(table)

	table.float = &canvas.Image{}
	table.float.Hide()
	return table
}
