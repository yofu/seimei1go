package seimei1go

type Event interface {
}

type Key int

const (
	KeyUnknown Key = iota
	KeyEsc
)

type EventKey struct {
	Key Key
}

type MouseState int

const (
	MouseUnknown MouseState = iota
	MouseRelease
)

type EventMouse struct {
	Key MouseState
	X   int
	Y   int
}
