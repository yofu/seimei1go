package termboxdriver

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/yofu/seimei1go"
)

func Key(k termbox.Key) seimei1go.Key {
	switch k {
	case termbox.KeyEsc:
		return seimei1go.KeyEsc
	default:
		return seimei1go.KeyUnknown
	}
}

func MouseState(s termbox.Key) seimei1go.MouseState {
	switch s {
	case termbox.MouseRelease:
		return seimei1go.MouseRelease
	default:
		return seimei1go.MouseUnknown
	}
}
