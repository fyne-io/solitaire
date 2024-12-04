package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fyne-io/solitaire/faces"
)

const minCardWidth = float32(95)
const minPadding = float32(4)
const cardRatio = 142.0 / minCardWidth

var (
	cardSize = fyne.Size{Width: minCardWidth, Height: minCardWidth * cardRatio}

	smallPad = minPadding
	overlap  = smallPad * 5
	bigPad   = smallPad + overlap

	minWidth  = cardSize.Width*7 + smallPad*6
	minHeight = cardSize.Height*3 + smallPad*2 + 1
)

func updateCardPosition(c *canvas.Image, x, y float32) {
	c.Resize(cardSize)
	c.Move(fyne.NewPos(x, y))
}

func withinCardBounds(c *canvas.Image, pos fyne.Position) bool {
	if pos.X < c.Position().X || pos.Y < c.Position().Y {
		return false
	}

	if pos.X >= c.Position().X+c.Size().Width || pos.Y >= c.Position().Y+c.Size().Height {
		return false
	}

	return true
}

func newCardPos(card *Card) *canvas.Image {
	if card == nil {
		return &canvas.Image{}
	}

	var face fyne.Resource
	if card.FaceUp {
		face = card.Face()
	} else {
		face = faces.ForBack()
	}
	image := &canvas.Image{Resource: face}
	image.Resize(cardSize)

	return image
}

func newCardSpace() *canvas.Image {
	space := faces.ForSpace()
	image := &canvas.Image{Resource: space}
	image.Resize(cardSize)

	return image
}

type tableRender struct {
	game *Game

	deck *canvas.Image
	sep  *widget.Separator

	pile1, pile2, pile3            *canvas.Image
	build1, build2, build3, build4 *canvas.Image

	stack1, stack2, stack3, stack4, stack5, stack6, stack7 *stackRender

	objects []fyne.CanvasObject
	table   *Table
}

func updateSizes(pad float32) {
	smallPad = pad
	overlap = smallPad * 5
}

func (t *tableRender) MinSize() fyne.Size {
	return fyne.NewSize(220, 140)
}

func (t *tableRender) Layout(size fyne.Size) {
	padding := size.Width * .006
	updateSizes(padding)

	newWidth := (size.Width - smallPad*6) / 7.0
	cardSize = fyne.NewSize(newWidth, newWidth*cardRatio)

	updateCardPosition(t.deck, 0, 0)

	updateCardPosition(t.pile1, smallPad+cardSize.Width, 0)
	updateCardPosition(t.pile2, smallPad+cardSize.Width+overlap, 0)
	updateCardPosition(t.pile3, smallPad+cardSize.Width+overlap*2, 0)
	updateCardPosition(t.table.float[0], 0, 0)

	updateCardPosition(t.build1, size.Width-smallPad*3-cardSize.Width*4, 0)
	updateCardPosition(t.build2, size.Width-smallPad*2-cardSize.Width*3, 0)
	updateCardPosition(t.build3, size.Width-smallPad-cardSize.Width*2, 0)
	updateCardPosition(t.build4, size.Width-cardSize.Width, 0)

	sepThick := theme.SeparatorThicknessSize()
	t.sep.Resize(fyne.NewSize(size.Width, sepThick))
	t.sep.Move(fyne.NewPos(0, cardSize.Height+smallPad))

	stackYPos := smallPad*2 + sepThick + cardSize.Height
	t.stack1.Layout(fyne.NewPos(0, stackYPos),
		fyne.NewSize(cardSize.Width, size.Height-stackYPos))
	t.stack2.Layout(fyne.NewPos(smallPad+cardSize.Width, stackYPos),
		fyne.NewSize(cardSize.Width, size.Height-stackYPos))
	t.stack3.Layout(fyne.NewPos((smallPad+cardSize.Width)*2, stackYPos),
		fyne.NewSize(cardSize.Width, size.Height-stackYPos))
	t.stack4.Layout(fyne.NewPos((smallPad+cardSize.Width)*3, stackYPos),
		fyne.NewSize(cardSize.Width, size.Height-stackYPos))
	t.stack5.Layout(fyne.NewPos((smallPad+cardSize.Width)*4, stackYPos),
		fyne.NewSize(cardSize.Width, size.Height-stackYPos))
	t.stack6.Layout(fyne.NewPos((smallPad+cardSize.Width)*5, stackYPos),
		fyne.NewSize(cardSize.Width, size.Height-stackYPos))
	t.stack7.Layout(fyne.NewPos((smallPad+cardSize.Width)*6, stackYPos),
		fyne.NewSize(cardSize.Width, size.Height-stackYPos))
}

func (t *tableRender) ApplyTheme() {
	// no-op we are a custom UI
}

func (t *tableRender) refreshCard(img *canvas.Image, card *Card) {
	img.Hidden = card == nil
	t.refreshCardOrBlank(img, card)
}

func (t *tableRender) refreshCardOrBlank(img *canvas.Image, card *Card) {
	img.Translucency = 0
	if card == nil {
		img.Resource = faces.ForSpace()
		return
	}

	if card.FaceUp {
		img.Resource = card.Face()
	} else {
		img.Resource = faces.ForBack()
	}

	if t.table.selected != nil && cardEquals(card, t.table.selected) {
		img.Translucency = 0.25
	} else {
		img.Translucency = 0
	}
	img.Refresh()
}

func (t *tableRender) Refresh() {
	if len(t.game.Hand.Cards) > 0 {
		t.deck.Resource = faces.ForBack()
	} else {
		t.deck.Resource = faces.ForSpace()
	}
	canvas.Refresh(t.deck)
	canvas.Refresh(t.sep)

	t.refreshCard(t.pile1, t.game.Draw1)
	t.refreshCard(t.pile2, t.game.Draw2)
	t.refreshCard(t.pile3, t.game.Draw3)

	t.refreshCardOrBlank(t.build1, t.game.Build1.Top())
	t.refreshCardOrBlank(t.build2, t.game.Build2.Top())
	t.refreshCardOrBlank(t.build3, t.game.Build3.Top())
	t.refreshCardOrBlank(t.build4, t.game.Build4.Top())

	t.stack1.Refresh(t.game.Stack1)
	t.stack2.Refresh(t.game.Stack2)
	t.stack3.Refresh(t.game.Stack3)
	t.stack4.Refresh(t.game.Stack4)
	t.stack5.Refresh(t.game.Stack5)
	t.stack6.Refresh(t.game.Stack6)
	t.stack7.Refresh(t.game.Stack7)

	canvas.Refresh(t.table)
}

func (t *tableRender) Objects() []fyne.CanvasObject {
	return t.objects
}

func (t *tableRender) Destroy() {
}

func (t *tableRender) appendStack(stack *stackRender) {
	for _, card := range stack.cards {
		t.objects = append(t.objects, card)
	}
}

func (t *tableRender) positionForCard(_ *Card) *canvas.Image {
	return nil
}

func (t *tableRender) findCard(pos fyne.Position) ([]*Card, []*canvas.Image, bool) {
	if t.game.Draw3 != nil {
		if withinCardBounds(t.pile3, pos) {
			return []*Card{t.game.Draw3}, []*canvas.Image{t.pile3}, false
		}
	} else if t.game.Draw2 != nil {
		if withinCardBounds(t.pile2, pos) {
			return []*Card{t.game.Draw2}, []*canvas.Image{t.pile2}, false
		}
	} else if t.game.Draw1 != nil {
		if withinCardBounds(t.pile1, pos) {
			return []*Card{t.game.Draw1}, []*canvas.Image{t.pile1}, true
		}
	}

	// Skipping build piles as we can't drag out...

	if c, p := t.findOnStack(t.stack1, t.game.Stack1, pos); c != nil {
		return c, p, len(t.game.Stack1.Cards) == 1
	} else if c, p := t.findOnStack(t.stack2, t.game.Stack2, pos); c != nil {
		return c, p, len(t.game.Stack2.Cards) == 1
	} else if c, p := t.findOnStack(t.stack3, t.game.Stack3, pos); c != nil {
		return c, p, len(t.game.Stack3.Cards) == 1
	} else if c, p := t.findOnStack(t.stack4, t.game.Stack4, pos); c != nil {
		return c, p, len(t.game.Stack4.Cards) == 1
	} else if c, p := t.findOnStack(t.stack5, t.game.Stack5, pos); c != nil {
		return c, p, len(t.game.Stack5.Cards) == 1
	} else if c, p := t.findOnStack(t.stack6, t.game.Stack6, pos); c != nil {
		return c, p, len(t.game.Stack6.Cards) == 1
	} else if c, p := t.findOnStack(t.stack7, t.game.Stack7, pos); c != nil {
		return c, p, len(t.game.Stack7.Cards) == 1
	}

	return nil, nil, false
}

func (t *tableRender) findOnStack(render *stackRender, stack *Stack, pos fyne.Position) ([]*Card, []*canvas.Image) {
	for i := len(stack.Cards) - 1; i >= 0; i-- {
		if withinCardBounds(render.cards[i], pos) {
			return stack.Cards[i:len(stack.Cards)],render.cards[i:len(stack.Cards)]
		}
	}

	return nil, nil
}

func (t *tableRender) stackPos(i int) fyne.Position {
	size := t.table.Size()
	switch i {
	case 3:
		return fyne.NewPos(size.Width-cardSize.Width, 0)
	case 2:
		return fyne.NewPos(size.Width-smallPad-cardSize.Width*2, 0)
	case 1:
		return fyne.NewPos(size.Width-smallPad*2-cardSize.Width*3, 0)
	default:
		return fyne.NewPos(size.Width-smallPad*3-cardSize.Width*4, 0)
	}
}

func newTableRender(table *Table) *tableRender {
	render := &tableRender{}
	table.findCard = render.findCard
	table.stackPos = render.stackPos
	render.table = table
	render.game = table.game
	render.deck = newCardPos(nil)
	render.sep = widget.NewSeparator()

	render.pile1 = newCardPos(nil)
	render.pile2 = newCardPos(nil)
	render.pile3 = newCardPos(nil)

	render.build1 = newCardSpace()
	render.build2 = newCardSpace()
	render.build3 = newCardSpace()
	render.build4 = newCardSpace()

	render.stack1 = newStackRender(render)
	render.stack2 = newStackRender(render)
	render.stack3 = newStackRender(render)
	render.stack4 = newStackRender(render)
	render.stack5 = newStackRender(render)
	render.stack6 = newStackRender(render)
	render.stack7 = newStackRender(render)

	render.objects = []fyne.CanvasObject{render.deck, render.sep, render.pile1, render.pile2, render.pile3,
		render.build1, render.build2, render.build3, render.build4}

	render.appendStack(render.stack1)
	render.appendStack(render.stack2)
	render.appendStack(render.stack3)
	render.appendStack(render.stack4)
	render.appendStack(render.stack5)
	render.appendStack(render.stack6)
	render.appendStack(render.stack7)

	floats := container.NewWithoutLayout()
	for i := 0; i < len(table.float); i++ {
		floats.Add(table.float[i])
	}
	render.objects = append(render.objects, floats)
	render.Refresh()
	return render
}

type stackRender struct {
	cards [19]*canvas.Image
	table *tableRender
}

func (s *stackRender) Layout(pos fyne.Position, _ fyne.Size) {
	top := pos.Y
	for i := range s.cards {
		updateCardPosition(s.cards[i], pos.X, top)

		top += overlap
	}
}

func (s *stackRender) Refresh(stack *Stack) {
	var i int
	var card *Card
	if len(stack.Cards) == 0 {
		s.cards[0].Resource = faces.ForSpace()
		s.cards[0].Translucency = 0
		i = 0
	} else {
		for i, card = range stack.Cards {
			s.table.refreshCard(s.cards[i], card)
		}
	}

	for i = i + 1; i < len(s.cards); i++ {
		s.cards[i].Image = nil
		s.cards[i].Resource = nil
		s.cards[i].Hide()
	}
}

func newStackRender(table *tableRender) *stackRender {
	r := &stackRender{table: table}
	for i := range r.cards {
		r.cards[i] = newCardPos(nil)
	}

	return r
}
