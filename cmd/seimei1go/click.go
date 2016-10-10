package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/google/subcommands"
	termbox "github.com/nsf/termbox-go"
	"github.com/yofu/seimei1go"
	"golang.org/x/net/context"
)

type click struct {
	N      int
	moving bool
}

func (*click) Name() string {
	return "click"
}

func (*click) Synopsis() string {
	return "create a hole at clicked point"
}

func (*click) Usage() string {
	return "click"
}

func (c *click) SetFlags(f *flag.FlagSet) {
	f.IntVar(&c.N, "size", 32, "board size")
}

func (c *click) draw(b *seimei1go.Board) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for i := 0; i < c.N; i++ {
		termbox.SetCell(i, 0, ' ', termbox.ColorDefault, termbox.ColorWhite)
		termbox.SetCell(i, c.N, ' ', termbox.ColorDefault, termbox.ColorWhite)
		termbox.SetCell(0, i, ' ', termbox.ColorDefault, termbox.ColorWhite)
		termbox.SetCell(c.N, i, ' ', termbox.ColorDefault, termbox.ColorWhite)
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

func (c *click) pollEvent(b *seimei1go.Board) {
	c.draw(b)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
			default:
				c.draw(b)
			}
		case termbox.EventMouse:
			x := ev.MouseX
			y := ev.MouseY
			if !c.moving && ev.Key == termbox.MouseRelease && b.State(x, y) == seimei1go.BLANK {
				h, err := b.CreateHole(x, y)
				if err != nil {
					continue
				}
				c.moving = true
				go func(board *seimei1go.Board, hole *seimei1go.Hole) {
					for {
						select {
						case <-time.After(20 * time.Millisecond):
							err := hole.Move()
							if err != nil {
								board.SetBound()
								c.draw(board)
								c.moving = false
								return
							}
							c.draw(board)
						}
					}
				}(b, h)
			}
		default:
			c.draw(b)
		}
	}
}

func (c *click) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	rand.Seed(time.Now().UnixNano())
	b := seimei1go.NewBoard(c.N, c.N)
	b.Set(4, 2, seimei1go.BOUND)
	b.Set(3, 3, seimei1go.BOUND)
	b.Set(4, 3, seimei1go.INNER)
	b.Set(5, 3, seimei1go.BOUND)
	b.Set(2, 4, seimei1go.BOUND)
	b.Set(3, 4, seimei1go.INNER)
	b.Set(4, 4, seimei1go.INNER)
	b.Set(5, 4, seimei1go.INNER)
	b.Set(6, 4, seimei1go.BOUND)
	b.Set(3, 5, seimei1go.BOUND)
	b.Set(4, 5, seimei1go.INNER)
	b.Set(5, 5, seimei1go.BOUND)
	b.Set(4, 6, seimei1go.BOUND)
	c.pollEvent(b)
	return subcommands.ExitSuccess
}
