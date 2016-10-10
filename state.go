package seimei1go

type State int

const (
	OUTOFRANGE State = iota
	BLANK
	BOUND
	INNER
)
