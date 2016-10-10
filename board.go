package seimei1go

import "fmt"

type Board struct {
	X    int
	Y    int
	data [][]*Point
}

func NewBoard(m, n int) *Board {
	d := make([][]*Point, m)
	for i := 0; i < m; i++ {
		d[i] = make([]*Point, n)
		for j := 0; j < n; j++ {
			d[i][j] = NewPoint(i, j)
		}
	}
	return &Board{
		X:    m,
		Y:    n,
		data: d,
	}
}

func (b *Board) Print() {
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

func (b *Board) SetBound() {
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

func (b *Board) State(x, y int) State {
	if x < 0 || x >= b.X || y < 0 || y >= b.Y {
		return OUTOFRANGE
	}
	return b.data[x][y].state
}

func (b *Board) Set(x, y int, s State) {
	if x < 0 || x >= b.X || y < 0 || y >= b.Y {
		return
	}
	b.data[x][y].state = s
}

func (b *Board) CreateHole(x, y int) (*Hole, error) {
	min := b.X + b.Y
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
