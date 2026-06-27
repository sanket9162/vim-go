package mode

import "github.com/gdamore/tcell/v3"

type NormalMode struct {
	pendingKey string
	count      int // Tacks nuberic respeat prefix (e.g 5 for 5j)
}

// Name returns the name of the normal mode.
func (m *NormalMode) Name() string { return "NORMAL" }

// HandleKey handles the key press events in normal mode.
func (m *NormalMode) HandleKey(e EditorInterface, ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyCtrlR:
		e.Redo()
		return
	case tcell.KeyEscape:
		m.pendingKey = ""
		m.count = 0 // Reset count
		e.QuitEditor()
	case tcell.KeyRune:
		keyStr := ev.Str()

		//Accumulate numeric prefix digits (e.g '5' in 5j)
		if keyStr >= "0" && keyStr <= "9" {
			// '0' is a jump to start of line unless it's part of cound (e.g. '10')
			if keyStr == "0" && m.count == 0 {
				// Fall through to normal command matching
			} else {
				digit := int(keyStr[0] - '0')
				m.count = m.count*10 + digit
				return
			}
		}

		// Calculate how many times to repeat the movement/command
		repeat := 1
		if m.count > 0 {
			repeat = m.count
		}
		if m.pendingKey == "d" {
			m.pendingKey = ""
			switch keyStr {
			case "d":
				// Support multi-line deletion (e.g 2dd)
				for i := 0; i < repeat; i++ {
					e.DeleteLine()
				}
				m.count = 0
				return
			case "w":
				// Support multi-word deletion (e.g 2dw)
				for i := 0; i < repeat; i++ {
					e.DeleteWord()
				}
				return
			}
			return
		}
		if m.pendingKey == "g" {
			m.pendingKey = ""
			if keyStr == "g" {
				e.MoveCursorToStartOfFile()
				return
			}
		}

		switch keyStr {
		case "d":
			m.pendingKey = "d"
		case "x":
			for i := 0; i < repeat; i++ {
				e.DeleteUnderCursor()
			}
			m.count = 0
		case "g":
			m.pendingKey = "g"
		case "h":
			for i := 0; i < repeat; i++ {
				e.MoveCursorLeft()
			}
			m.count = 0
		case "l":
			for i := 0; i < repeat; i++ {
				e.MoveCursorRight()
			}
			m.count = 0
		case "j":
			for i := 0; i < repeat; i++ {
				e.MoveCursorDown()
			}
			m.count = 0
		case "k":
			for i := 0; i < repeat; i++ {
				e.MoveCursorUp()
			}
			m.count = 0
		case "0":
			e.MoveCursorToStartOfLine()
			m.count = 0
		case "G":
			e.MoveCursorToEndOfFile()
			m.count = 0
		case "w":
			for i := 0; i < repeat; i++ {
				e.MoveCursorToNextWord()
			}
			m.count = 0
		case "a":
			e.MoveCursorRight()
			e.SetMode("INSERT")
			m.count = 0
		case "A":
			e.MoveCursorToEndOfLine()
			e.SetMode("INSERT")
			m.count = 0
		case "$":
			e.MoveCursorToEndOfLine()
			m.count = 0
		case "p":
			e.Paste(false)
			m.count = 0
		case "P":
			e.Paste(true)
			m.count = 0
		case "u":
			for i := 0; i < repeat; i++ {
				e.Undo()
			}
			m.count = 0
		case "n":
			for i := 0; i < repeat; i++ {
				e.SearchNext()
			}
			m.count = 0
		case "N":
			for i := 0; i < repeat; i++ {
				e.SearchPrev()
			}
			m.count = 0
		case "i":
			e.SetMode("INSERT")
			m.count = 0
		case "v":
			e.SetMode("VISUAL")
			m.count = 0
		case ":":
			e.SetMode("COMMAND")
			m.count = 0
		case "/":
			e.SetMode("SEARCH")
			m.count = 0
		default:
			// Reset count on any other key
			m.count = 0
		}
	default:
		m.pendingKey = ""
		m.count = 0
	}
}
