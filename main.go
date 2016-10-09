package main

import (
	"fmt"
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

const (
	N = 32
)

var (
	moving = false
)

type state int

const (
	OUTOFRANGE state = iota
	BLANK
	BOUND
	INNER
)

type point struct {
	X     int
	Y     int
	state state
}

func NewPoint(x, y int) *point {
	return &point{
		X:     x,
		Y:     y,
		state: BLANK,
	}
}

type board struct {
	X    int
	Y    int
	data [][]*point
}

func NewBoard(m, n int) *board {
	d := make([][]*point, m)
	for i := 0; i < m; i++ {
		d[i] = make([]*point, n)
		for j := 0; j < n; j++ {
			d[i][j] = NewPoint(i, j)
		}
	}
	return &board{
		X:    m,
		Y:    n,
		data: d,
	}
}

func (b *board) Print() {
	for i := 0; i < b.Y; i++ {
		for j := 0; j < b.X; j++ {
			switch b.data[j][i].state {
			case BLANK:
				fmt.Print("_")
			case BOUND:
				fmt.Print("*")
			case INNER:
				fmt.Print("+")
			}
		}
		fmt.Println("")
	}
}

func (b *board) SetBound() {
	for i := 0; i < b.X; i++ {
		for j := 0; j < b.Y; j++ {
			if b.data[i][j].state == BLANK {
				continue
			}
			if s := b.State(i-1, j); s == OUTOFRANGE || s == BLANK {
				b.data[i][j].state = BOUND
				continue
			}
			if s := b.State(i+1, j); s == OUTOFRANGE || s == BLANK {
				b.data[i][j].state = BOUND
				continue
			}
			if s := b.State(i, j-1); s == OUTOFRANGE || s == BLANK {
				b.data[i][j].state = BOUND
				continue
			}
			if s := b.State(i, j+1); s == OUTOFRANGE || s == BLANK {
				b.data[i][j].state = BOUND
				continue
			}
			b.data[i][j].state = INNER
		}
	}
}

func (b *board) State(x, y int) state {
	if x < 0 || x >= b.X || y < 0 || y >= b.Y {
		return OUTOFRANGE
	}
	return b.data[x][y].state
}

func (b *board) Set(x, y int, s state) {
	if x < 0 || x >= b.X || y < 0 || y >= b.Y {
		return
	}
	b.data[x][y].state = s
}

func (b *board) CreateHole(x, y int) (*hole, error) {
	min := 2 * N
	var hx, hy int
	create := false
	abs := func(val int) int {
		if val < 0 {
			return -val
		}
		return val
	}
	for i := 0; i < b.X; i++ {
		for j := 0; j < b.Y; j++ {
			if b.data[i][j].state == BOUND {
				dx := abs(i - x)
				dy := abs(j - y)
				d := dx + dy
				if d < min {
					create = true
					if dx < dy {
						if j < y {
							hx = i
							hy = j + 1
						} else {
							hx = i
							hy = j - 1
						}
					} else {
						if i < x {
							hx = i + 1
							hy = j
						} else {
							hx = i - 1
							hy = j
						}
					}
					min = d
				}
			}
		}
	}
	if !create || b.State(hx, hy) != BLANK {
		return nil, fmt.Errorf("cannot create hole")
	} else {
		return NewHole(b, hx, hy), nil
	}
}

type hole struct {
	board    *board
	position *point
	memory   []*point
}

func NewHole(b *board, x, y int) *hole {
	return &hole{
		board: b,
		position: &point{
			X:     x,
			Y:     y,
			state: BLANK,
		},
		memory: []*point{NewPoint(x, y)},
	}
}

func (h *hole) Up() *point {
	y := h.position.Y + 1
	if y >= h.board.Y {
		y -= h.board.Y
	}
	return h.board.data[h.position.X][y]
}

func (h *hole) Down() *point {
	y := h.position.Y - 1
	if y < 0 {
		y += h.board.Y
	}
	return h.board.data[h.position.X][y]
}

func (h *hole) Right() *point {
	x := h.position.X + 1
	if x >= h.board.X {
		x -= h.board.X
	}
	return h.board.data[x][h.position.Y]
}

func (h *hole) Left() *point {
	x := h.position.X - 1
	if x < 0 {
		x += h.board.X
	}
	return h.board.data[x][h.position.Y]
}

func (h *hole) Move() error {
	cand := make([]*point, 0)
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

func draw(b *board) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for i := 0; i < N; i++ {
		termbox.SetCell(i, 0, ' ', termbox.ColorDefault, termbox.ColorWhite)
		termbox.SetCell(i, N, ' ', termbox.ColorDefault, termbox.ColorWhite)
		termbox.SetCell(0, i, ' ', termbox.ColorDefault, termbox.ColorWhite)
		termbox.SetCell(N, i, ' ', termbox.ColorDefault, termbox.ColorWhite)
	}
	for i := 0; i < b.X; i++ {
		for j := 0; j < b.Y; j++ {
			var color termbox.Attribute
			switch b.data[i][j].state {
			case BLANK:
				color = termbox.ColorDefault
			case INNER:
				color = termbox.ColorYellow
			case BOUND:
				color = termbox.ColorRed
			}
			termbox.SetCell(i, j, ' ', termbox.ColorDefault, color)
		}
	}
	termbox.Flush()
}

func pollEvent(b *board) {
	draw(b)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
			default:
				draw(b)
			}
		case termbox.EventMouse:
			x := ev.MouseX
			y := ev.MouseY
			if !moving && ev.Key == termbox.MouseRelease && b.State(x, y) == BLANK {
				h, err := b.CreateHole(x, y)
				if err != nil {
					continue
				}
				moving = true
				go func(board *board, hole *hole) {
					for {
						select {
						case <-time.After(20 * time.Millisecond):
							err := hole.Move()
							if err != nil {
								board.SetBound()
								draw(board)
								moving = false
								return
							}
							draw(board)
						}
					}
				}(b, h)
			}
		default:
			draw(b)
		}
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	rand.Seed(time.Now().UnixNano())
	b := NewBoard(N, N)
	b.Set(4, 2, BOUND)
	b.Set(3, 3, BOUND)
	b.Set(4, 3, INNER)
	b.Set(5, 3, BOUND)
	b.Set(2, 4, BOUND)
	b.Set(3, 4, INNER)
	b.Set(4, 4, INNER)
	b.Set(5, 4, INNER)
	b.Set(6, 4, BOUND)
	b.Set(3, 5, BOUND)
	b.Set(4, 5, INNER)
	b.Set(5, 5, BOUND)
	b.Set(4, 6, BOUND)
	pollEvent(b)
}
