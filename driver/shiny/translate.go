package shinydriver

import (
	"github.com/yofu/seimei1go"
	"golang.org/x/mobile/event/key"
)

func Key(c key.Code) seimei1go.Key {
	switch c {
	case key.CodeEscape:
		return seimei1go.KeyEsc
	default:
		return seimei1go.KeyUnknown
	}
}
