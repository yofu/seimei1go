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

type light struct {
	N      int
	X      int
	Y      int
	moving bool
}

func (*light) Name() string {
	return "light"
}

func (*light) Synopsis() string {
	return "Phototaxis simulation"
}

func (*light) Usage() string {
	return "light X Y"
}

func (c *light) SetFlags(f *flag.FlagSet) {
	f.IntVar(&c.N, "size", 64, "board size")
	f.IntVar(&c.X, "X", 60, "X coord")
	f.IntVar(&c.Y, "Y", 55, "Y coord")
}

func (l *light) draw(b *seimei1go.Board) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for i := 0; i < l.N; i++ {
		termbox.SetCell(i, 0, ' ', termbox.ColorDefault, termbox.ColorWhite)
		termbox.SetCell(i, l.N, ' ', termbox.ColorDefault, termbox.ColorWhite)
		termbox.SetCell(0, i, ' ', termbox.ColorDefault, termbox.ColorWhite)
		termbox.SetCell(l.N, i, ' ', termbox.ColorDefault, termbox.ColorWhite)
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

func (l *light) pollEvent(b *seimei1go.Board) {
	l.draw(b)
	go func(b0 *seimei1go.Board) {
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				if !l.moving && b.State(l.X, l.Y) == seimei1go.BLANK {
					h, err := b.CreateHole(l.X, l.Y)
					if err != nil {
						continue
					}
					l.moving = true
					go func(board *seimei1go.Board, hole *seimei1go.Hole) {
						for {
							select {
							case <-time.After(time.Millisecond):
								err := hole.Move()
								if err != nil {
									board.SetBound()
									l.draw(board)
									l.moving = false
									return
								}
								l.draw(board)
							}
						}
					}(b0, h)
				}
			}
		}
	}(b)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
			default:
				l.draw(b)
			}
		default:
			l.draw(b)
		}
	}
}

func (l *light) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	rand.Seed(time.Now().UnixNano())
	b := seimei1go.NewBoard(l.N, l.N)
	for i := 4; i < 12; i++ {
		for j := 4; j < 12; j++ {
			b.Set(i, j, seimei1go.INNER)
		}
	}
	b.SetBound()
	l.pollEvent(b)
	return subcommands.ExitSuccess
}
