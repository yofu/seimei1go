package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/google/subcommands"
	"github.com/yofu/seimei1go"
	driver "github.com/yofu/seimei1go/driver/shiny"
	"golang.org/x/net/context"
)

type random struct {
	N      int
	moving bool
}

func (*random) Name() string {
	return "random"
}

func (*random) Synopsis() string {
	return "random"
}

func (*random) Usage() string {
	return "random X Y"
}

func (c *random) SetFlags(f *flag.FlagSet) {
	f.IntVar(&c.N, "size", 64, "board size")
}

func (l *random) pollEvent(b *seimei1go.Board, ch chan seimei1go.Event) {
	rand.Seed(time.Now().UnixNano())
	driver.Draw(b)
	go func(b0 *seimei1go.Board) {
		for {
			select {
			case <-time.After(time.Millisecond):
				if !l.moving {
					h, err := b0.MoveFromRandomBound(func(x, y int) float64 { return 1.0 })
					if err != nil {
						continue
					}
					l.moving = true
					go func(board *seimei1go.Board, hole *seimei1go.Hole) {
						for {
							err := hole.Move()
							if err != nil {
								board.SetBound()
								driver.Draw(board)
								l.moving = false
								return
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

func (l *random) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	b := seimei1go.NewBoard(l.N, l.N)
	start := int(0.25 * float64(l.N))
	end := int(0.75 * float64(l.N))
	for i := start; i < end; i++ {
		for j := start; j < end; j++ {
			b.Set(i, j, seimei1go.INNER)
		}
	}
	b.SetBound()
	driver.Start(b, l.pollEvent)
	return subcommands.ExitSuccess
}
