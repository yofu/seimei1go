package seimei1go

type Point struct {
	X     int
	Y     int
	state State
}

func NewPoint(x, y int) *Point {
	return &Point{
		X:     x,
		Y:     y,
		state: BLANK,
	}
}
