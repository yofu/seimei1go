package main

import (
	"flag"
	"time"

	"github.com/google/subcommands"
	"github.com/yofu/seimei1go"
	driver "github.com/yofu/seimei1go/driver/shiny"
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

func (l *light) pollEvent(b *seimei1go.Board, ch chan seimei1go.Event) {
	driver.Draw(b)
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
									driver.Draw(board)
									l.moving = false
									return
								}
								driver.Draw(board)
							}
						}
					}(b0, h)
				}
			}
		}
	}(b)
	for {
		select {
		case e := <-ch:
			switch ev := e.(type) {
			case seimei1go.EventKey:
				switch ev.Key {
				case seimei1go.KeyEsc:
					return
				}
			}
		}
	}
}

func (l *light) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	b := seimei1go.NewBoard(l.N, l.N)
	for i := 4; i < 12; i++ {
		for j := 4; j < 12; j++ {
			b.Set(i, j, seimei1go.INNER)
		}
	}
	b.SetBound()
	driver.Start(b, l.pollEvent)
	return subcommands.ExitSuccess
}
