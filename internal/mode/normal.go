package mode

import "github.com/gdamore/tcell/v3"

type NormalMode struct {
	pendingKey string
}

// Name returns the name of the normal mode.
func (m *NormalMode) Name() string { return "NORMAL" }

// HandleKey handles the key press events in normal mode.
func (m *NormalMode) HandleKey(e EditorInterface, ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		m.pendingKey = ""
		e.QuitEditor()
	case tcell.KeyRune:
		keyStr := ev.Str()
		if m.pendingKey == "g" {
			m.pendingKey = ""
			if keyStr == "g" {
				e.MoveCursorToStartOfFile()
				return
			}
		}

		switch keyStr {
		case "g":
			m.pendingKey = "g"
		case "j":
			e.MoveCursorLeft()
		case ";":
			e.MoveCursorRight()
		case "k":
			e.MoveCursorDown()
		case "l":
			e.MoveCursorUp()
		case "0":
			e.MoveCursorToStartOfLine()
		case "G":
			e.MoveCursorToEndOfFile()
		case "w":
			e.MoveCursorToNextWord()
		case "a":
			e.MoveCursorRight()
			e.SetMode("INSERT")
		case "A":
			e.MoveCursorToEndOfLine()
			e.SetMode("INSERT")
		case "$":
			e.MoveCursorToEndOfLine()
		case "i":
			e.SetMode("INSERT")
		case ":":
			e.SetMode("COMMAND")
		}
	default:
		m.pendingKey = ""
	}
}
