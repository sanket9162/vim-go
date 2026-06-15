package mode

import "github.com/gdamore/tcell/v3"

type NormalMode struct{}

func (m *NormalMode) Name() string { return "NORMAL" }

func (m *NormalMode) HandleKey(e EditorInterface, ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		e.QuitEditor()
	case tcell.KeyRune:
		switch ev.Str() {
		case "h":
			e.MoveCursorLeft()
		case "l":
			e.MoveCursorRight()
		case "j":
			e.MoveCursorDown()
		case "k":
			e.MoveCursorUp()
		case "i":
			e.SetMode("INSERT")
		}
	}
}
