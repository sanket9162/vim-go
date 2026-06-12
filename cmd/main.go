package main

import (
	"log"

	"github.com/gdamore/tcell/v3"
)

func main() {
	// initialize the screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+V", err)
	}
	if err = s.Init(); err != nil {
		log.Fatalf("%+V", err)
	}
	defer s.Fini()

	// Default style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	// Cursor state
	cx, cy := 0, 0
	mode := "NORMAL"

	for {
		// 1. Draw UI
		s.Clear()

		// Draw mode indicator at the bottom
		_, h := s.Size()
		statusLine := "-- " + mode + " --"
		drawText(s, 0, h-1, defStyle, statusLine)

		// Draw instructions
		drawText(s, 0, 0, defStyle, "Vim-Go Prototype")
		drawText(s, 0, 1, defStyle, "h,j,k,l: move | i: insert | Esc: exit (NORMAL)")

		// Set hardware cursor position
		s.ShowCursor(cx, cy)
		s.Show()

		// 2. Handle Events
		ev := <-s.EventQ()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			// Global exit
			if ev.Key() == tcell.KeyCtrlC {
				return
			}

			// Handle keys based on mode
			if mode == "NORMAL" {
				switch ev.Key() {
				case tcell.KeyEscape:
					return // Exit on Esc in Normal mode for now
				case tcell.KeyRune:
					switch ev.Str() {
					case "h":
						if cx > 0 {
							cx--
						}
					case "l":
						cx++
					case "j":
						cy++
					case "k":
						if cy > 0 {
							cy--
						}
					case "i":
						mode = "INSERT"
					}
				}
			} else if mode == "INSERT" {
				if ev.Key() == tcell.KeyEscape {
					mode = "NORMAL"
					if cx > 0 {
						cx--
					}
				}
				// In insert mode, we just move the cursor for now to show it's working
				if ev.Key() == tcell.KeyRune {
					cx++
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
