package main

import (
	"log"

	"github.com/gdamore/tcell/v3"
	"github.com/sanket9162/vim-go/internal/buffer"
)

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+V", err)
	}
	if err = s.Init(); err != nil {
		log.Fatalf("%+V", err)
	}
	defer s.Fini()

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	buffer := buffer.NewBuffer()
	cx, cy := 0, 0
	mode := "NORMAL"

	for {
		s.Clear()

		// 1. Render Buffer
		for y, line := range buffer.Lines {
			for x, r := range line {
				s.SetContent(x, y, r, nil, defStyle)
			}
		}

		// 2. Render Status Line
		_, h := s.Size()
		statusLine := "-- " + mode + " --"
		drawText(s, 0, h-1, defStyle, statusLine)

		s.ShowCursor(cx, cy)
		s.Show()

		// 3. Handle Events
		ev := <-s.EventQ()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				return
			}

			if mode == "NORMAL" {
				switch ev.Key() {
				case tcell.KeyEscape:
					return
				case tcell.KeyRune:
					switch ev.Str() {
					case "h":
						if cx > 0 {
							cx--
						}
					case "l":
						if cx < len(buffer.Lines[cy]) {
							cx++
						}
					case "j":
						if cy < len(buffer.Lines)-1 {
							cy++
							if cx > len(buffer.Lines[cy]) {
								cx = len(buffer.Lines[cy])
							}
						}
					case "k":
						if cy > 0 {
							cy--
							if cx > len(buffer.Lines[cy]) {
								cx = len(buffer.Lines[cy])
							}
						}
					case "i":
						mode = "INSERT"
					}
				}
			} else if mode == "INSERT" {
				switch ev.Key() {
				case tcell.KeyEscape:
					mode = "NORMAL"
					if cx > 0 {
						cx--
					}
				case tcell.KeyEnter:
					buffer.InsertNewline(cy, cx)
					cy++
					cx = 0
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					cy, cx = buffer.DeleteChar(cy, cx)
				case tcell.KeyRune:
					runes := []rune(ev.Str())
					for _, r := range runes {
						buffer.InsertChar(cy, cx, r)
						cx++
					}
				}
			}
		}
	}
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}
