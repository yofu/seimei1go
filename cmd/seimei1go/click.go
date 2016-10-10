package main

import (
	"flag"
	"time"

	"github.com/google/subcommands"
	"github.com/yofu/seimei1go"
	driver "github.com/yofu/seimei1go/driver/shiny"
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

func (c *click) pollEvent(b *seimei1go.Board, ch chan seimei1go.Event) {
	driver.Draw(b)
	for {
		select {
		case e := <-ch:
			switch ev := e.(type) {
			case seimei1go.EventKey:
				switch ev.Key {
				case seimei1go.KeyEsc:
					return
				}
			case seimei1go.EventMouse:
				x := ev.X
				y := ev.Y
				if !c.moving && ev.Key == seimei1go.MouseRelease && b.State(x, y) == seimei1go.BLANK {
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
									driver.Draw(board)
									c.moving = false
									return
								}
								driver.Draw(board)
							}
						}
					}(b, h)
				}
			}
		}
	}
}

func (c *click) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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
	driver.Start(b, c.pollEvent)
	return subcommands.ExitSuccess
}
