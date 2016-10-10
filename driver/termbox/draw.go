package termboxdriver

import (
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/yofu/seimei1go"
)

func Draw(b *seimei1go.Board) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for i := 0; i < b.X; i++ {
		termbox.SetCell(i, 0, ' ', termbox.ColorDefault, termbox.ColorWhite)
		termbox.SetCell(i, b.Y, ' ', termbox.ColorDefault, termbox.ColorWhite)
	}
	for i := 0; i < b.Y; i++ {
		termbox.SetCell(0, i, ' ', termbox.ColorDefault, termbox.ColorWhite)
		termbox.SetCell(b.X, i, ' ', termbox.ColorDefault, termbox.ColorWhite)
	}
	for i := 0; i < b.X; i++ {
		for j := 0; j < b.Y; j++ {
			var color termbox.Attribute
			switch b.State(i, j) {
			case seimei1go.BLANK:
				color = termbox.ColorDefault
			case seimei1go.INNER:
				color = termbox.ColorYellow
			case seimei1go.BOUND:
				color = termbox.ColorRed
			}
			termbox.SetCell(i, j, ' ', termbox.ColorDefault, color)
		}
	}
	termbox.Flush()
}

func Start(b *seimei1go.Board, poll func(*seimei1go.Board, chan seimei1go.Event)) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	rand.Seed(time.Now().UnixNano())
	ch := make(chan seimei1go.Event)
	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				ch <- seimei1go.EventKey{
					Key: Key(ev.Key),
				}
			case termbox.EventMouse:
				ch <- seimei1go.EventMouse{
					Key: MouseState(ev.Key),
					X:   ev.MouseX,
					Y:   ev.MouseY,
				}
			default:
			}
		}
	}()
	poll(b, ch)
}
