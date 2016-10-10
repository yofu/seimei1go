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

type feed struct {
	N      int
	moving bool
}

func (*feed) Name() string {
	return "feed"
}

func (*feed) Synopsis() string {
	return "feeding simulation"
}

func (*feed) Usage() string {
	return "feed"
}

func (c *feed) SetFlags(f *flag.FlagSet) {
	f.IntVar(&c.N, "size", 64, "board size")
}

func (l *feed) draw(b *seimei1go.Board) {
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

func (l *feed) pollEvent(b *seimei1go.Board) {
	l.draw(b)
	go func(b0 *seimei1go.Board) {
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				if !l.moving {
					h, err := b0.MoveFromRandomBound(func(x, y int) float64 {
						if x+y > int(0.6*float64(b0.X+b0.Y)) || x+y < int(0.4*float64(b0.X+b0.Y)) {
							return 0.7
						}
						return 0.3
					})
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

func (l *feed) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	rand.Seed(time.Now().UnixNano())
	b := seimei1go.NewBoard(l.N, l.N)
	for i := 22; i < 42; i++ {
		for j := 22; j < 42; j++ {
			b.Set(i, j, seimei1go.INNER)
		}
	}
	b.SetBound()
	l.pollEvent(b)
	return subcommands.ExitSuccess
}