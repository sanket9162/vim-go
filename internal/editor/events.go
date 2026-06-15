package editor

import (
	"github.com/gdamore/tcell/v3"
)

func (e *Editor) Run() {
	for {
		if e.Quit {
			return
		}

		e.Render()

		ev := <-e.Screen.EventQ()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			e.Screen.Sync()
		case *tcell.EventKey:
			e.handleKey(ev)
		}
	}
}

func (e *Editor) handleKey(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyCtrlC {
		e.Quit = true
		return
	}

	if e.Mode == "NORMAL" {
		e.handleNormalMode(ev)
	} else if e.Mode == "INSERT" {
		e.handleInsertMode(ev)
	}
}

func (e *Editor) handleNormalMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		e.Quit = true
	case tcell.KeyRune:
		switch ev.Str() {
		case "h":
			e.Cursor.MoveLeft()
		case "l":
			e.Cursor.MoveRight()
		case "j":
			e.Cursor.MoveDown()
		case "k":
			e.Cursor.MoveUp()
		case "i":
			e.Mode = "INSERT"
		}
	}
}

func (e *Editor) handleInsertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		e.Mode = "NORMAL"
		e.Cursor.MoveLeft()
	case tcell.KeyEnter:
		e.Buffer.InsertNewline(e.Cursor.Row(), e.Cursor.Col())
		e.Cursor.SetPos(0, e.Cursor.Row()+1)
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		row, col := e.Buffer.DeleteChar(e.Cursor.Row(), e.Cursor.Col())
		e.Cursor.SetPos(col, row)
	case tcell.KeyRune:
		runes := []rune(ev.Str())
		for _, r := range runes {
			e.Buffer.InsertChar(e.Cursor.Row(), e.Cursor.Col(), r)
			e.Cursor.SetPos(e.Cursor.Col()+1, e.Cursor.Row())
		}
	}
}
