package shinydriver

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/yofu/seimei1go"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var window screen.Window

const (
	pause = false
	play  = true
)

var pauseChan = make(chan bool, 64)

type uploadEvent struct{}

func Draw(b *seimei1go.Board) {
	if window != nil {
		window.Send(uploadEvent{})
	}
}

func drawpix(b *seimei1go.Board) []byte {
	rtn := make([]byte, 4*b.X*b.Y)
	for y := 0; y < b.Y; y++ {
		for x := 0; x < b.X; x++ {
			var v uint8
			switch b.State(x, y) {
			case seimei1go.BLANK:
				v = 0xff
			case seimei1go.BOUND:
				v = 0x00
			case seimei1go.INNER:
				v = 0x33
			}
			p := (b.X*y + x) * 4
			rtn[p+0] = v
			rtn[p+1] = v
			rtn[p+2] = v
			rtn[p+3] = 0xff
		}
	}
	return rtn
}

func Start(b *seimei1go.Board, poll func(*seimei1go.Board, chan seimei1go.Event)) {
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(nil)
		if err != nil {
			log.Fatal(err)
		}
		window = w
		buf, tex := screen.Buffer(nil), screen.Texture(nil)
		defer func() {
			if buf != nil {
				tex.Release()
				buf.Release()
			}
			w.Release()
		}()

		ch := make(chan seimei1go.Event)
		go poll(b, ch)
		var (
			sz size.Event
		)
		for {
			publish := false
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					pauseChan <- play
					var err error
					buf, err = s.NewBuffer(image.Point{b.X, b.Y})
					if err != nil {
						log.Fatal(err)
					}
					tex, err = s.NewTexture(image.Point{b.X, b.Y})
					if err != nil {
						log.Fatal(err)
					}
					tex.Fill(tex.Bounds(), color.White, draw.Src)
				case lifecycle.CrossOff:
					pauseChan <- pause
					tex.Release()
					tex = nil
					buf.Release()
					buf = nil
				}
			case key.Event:
				switch e.Direction {
				case key.DirPress:
					ch <- Key(e.Code)
				}
				if e.Code == key.CodeEscape {
					return
				}
			case mouse.Event:
				if e.Direction == mouse.DirRelease {
					ch <- seimei1go.EventMouse{
						Key: seimei1go.MouseRelease,
						X:   int(e.X),
						Y:   int(e.Y),
					}
				}
			case paint.Event:
				publish = buf != nil
			case size.Event:
				sz = e
			case uploadEvent:
				if buf != nil {
					pix := drawpix(b)
					copy(buf.RGBA().Pix, pix)
					publish = true
				}
				if publish {
					tex.Upload(image.Point{}, buf, buf.Bounds())
				}
			case error:
				log.Print(e)
			}
			if publish {
				w.Scale(sz.Bounds(), tex, tex.Bounds(), draw.Src, nil)
				w.Publish()
			}
		}
	})
}
