package seimei1go

import (
	"fmt"
	"math/rand"
)

type Hole struct {
	board    *Board
	position *Point
	memory   []*Point
}

func NewHole(b *Board, x, y int) *Hole {
	return &Hole{
		board: b,
		position: &Point{
			X:     x,
			Y:     y,
			state: BLANK,
		},
		memory: []*Point{NewPoint(x, y)},
	}
}

func (h *Hole) Up() *Point {
	y := h.position.Y + 1
	if y >= h.board.Y {
		y -= h.board.Y
	}
	return h.board.data[h.position.X][y]
}

func (h *Hole) Down() *Point {
	y := h.position.Y - 1
	if y < 0 {
		y += h.board.Y
	}
	return h.board.data[h.position.X][y]
}

func (h *Hole) Right() *Point {
	x := h.position.X + 1
	if x >= h.board.X {
		x -= h.board.X
	}
	return h.board.data[x][h.position.Y]
}

func (h *Hole) Left() *Point {
	x := h.position.X - 1
	if x < 0 {
		x += h.board.X
	}
	return h.board.data[x][h.position.Y]
}

func (h *Hole) Move() error {
	cand := make([]*Point, 0)
	num := 0
	if p := h.Up(); p.state != BLANK {
		add := true
		for _, m := range h.memory {
			if p.X == m.X && p.Y == m.Y {
				add = false
				continue
			}
		}
		if add {
			cand = append(cand, p)
			num++
		}
	}
	if p := h.Down(); p.state != BLANK {
		add := true
		for _, m := range h.memory {
			if p.X == m.X && p.Y == m.Y {
				add = false
				continue
			}
		}
		if add {
			cand = append(cand, p)
			num++
		}
	}
	if p := h.Right(); p.state != BLANK {
		add := true
		for _, m := range h.memory {
			if p.X == m.X && p.Y == m.Y {
				add = false
				continue
			}
		}
		if add {
			cand = append(cand, p)
			num++
		}
	}
	if p := h.Left(); p.state != BLANK {
		add := true
		for _, m := range h.memory {
			if p.X == m.X && p.Y == m.Y {
				add = false
				continue
			}
		}
		if add {
			cand = append(cand, p)
			num++
		}
	}
	if len(cand) == 0 {
		return fmt.Errorf("no candidate")
	}
	next := cand[rand.Int()%num]
	h.board.Set(h.position.X, h.position.Y, next.state)
	next.state = BLANK
	h.position.X = next.X
	h.position.Y = next.Y
	h.memory = append(h.memory, next)
	return nil
}

func (h *Hole) MoveUp() error {
	p := h.Up()
	if p.state == OUTOFRANGE {
		return fmt.Errorf("up is OUTOFRANGE")
	} else if p.state == BLANK {
		return fmt.Errorf("up is BLANK")
	}
	for _, m := range h.memory {
		if p.X == m.X && p.Y == m.Y {
			return fmt.Errorf("already visited")
		}
	}
	h.board.Set(h.position.X, h.position.Y, p.state)
	p.state = BLANK
	h.position.X = p.X
	h.position.Y = p.Y
	h.memory = append(h.memory, p)
	return nil
}

func (h *Hole) MoveDown() error {
	p := h.Down()
	if p.state == OUTOFRANGE {
		return fmt.Errorf("down is OUTOFRANGE")
	} else if p.state == BLANK {
		return fmt.Errorf("down is BLANK")
	}
	for _, m := range h.memory {
		if p.X == m.X && p.Y == m.Y {
			return fmt.Errorf("already visited")
		}
	}
	h.board.Set(h.position.X, h.position.Y, p.state)
	p.state = BLANK
	h.position.X = p.X
	h.position.Y = p.Y
	h.memory = append(h.memory, p)
	return nil
}

func (h *Hole) MoveRight() error {
	p := h.Right()
	if p.state == OUTOFRANGE {
		return fmt.Errorf("right is OUTOFRANGE")
	} else if p.state == BLANK {
		return fmt.Errorf("right is BLANK")
	}
	for _, m := range h.memory {
		if p.X == m.X && p.Y == m.Y {
			return fmt.Errorf("already visited")
		}
	}
	h.board.Set(h.position.X, h.position.Y, p.state)
	p.state = BLANK
	h.position.X = p.X
	h.position.Y = p.Y
	h.memory = append(h.memory, p)
	return nil
}

func (h *Hole) MoveLeft() error {
	p := h.Left()
	if p.state == OUTOFRANGE {
		return fmt.Errorf("left is OUTOFRANGE")
	} else if p.state == BLANK {
		return fmt.Errorf("left is BLANK")
	}
	for _, m := range h.memory {
		if p.X == m.X && p.Y == m.Y {
			return fmt.Errorf("already visited")
		}
	}
	h.board.Set(h.position.X, h.position.Y, p.state)
	p.state = BLANK
	h.position.X = p.X
	h.position.Y = p.Y
	h.memory = append(h.memory, p)
	return nil
}
