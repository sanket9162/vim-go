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

	b := buffer.NewBuffer()
	cursor := buffer.NewCursor(b)
	mode := "NORMAL"

	for {
		s.Clear()

		// 1. Render Buffer
		for y, line := range b.Lines {
			for x, r := range line {
				s.SetContent(x, y, r, nil, defStyle)
			}
		}

		// 2. Render Status Line
		_, h := s.Size()
		statusLine := "-- " + mode + " --"
		drawText(s, 0, h-1, defStyle, statusLine)

		s.ShowCursor(cursor.Col(), cursor.Row())
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
						cursor.MoveLeft()
					case "l":
						cursor.MoveRight()
					case "j":
						cursor.MoveDown()
					case "k":
						cursor.MoveUp()
					case "i":
						mode = "INSERT"
					}
				}
			} else if mode == "INSERT" {
				switch ev.Key() {
				case tcell.KeyEscape:
					mode = "NORMAL"
					cursor.MoveLeft()
				case tcell.KeyEnter:
					b.InsertNewline(cursor.Row(), cursor.Col())
					cursor.SetPos(0, cursor.Row()+1)
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					row, col := b.DeleteChar(cursor.Row(), cursor.Col())
					cursor.SetPos(col, row)
				case tcell.KeyRune:
					runes := []rune(ev.Str())
					for _, r := range runes {
						b.InsertChar(cursor.Row(), cursor.Col(), r)
						cursor.SetPos(cursor.Col()+1, cursor.Row())
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
