package mode

import "github.com/gdamore/tcell/v3"

type InsertMode struct {
	lastChar rune // Tracks the last character inserted to detecd sequences
}

// Name returns the name of the insert mode.
func (m *InsertMode) Name() string { return "INSERT" }

// HandleKey handles the key press events in insert mode.
func (m *InsertMode) HandleKey(e EditorInterface, ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		m.lastChar = 0
		e.SetMode("NORMAL")
	case tcell.KeyEnter:
		m.lastChar = 0
		e.InsertNewline()
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		m.lastChar = 0
		e.DeleteChar()
	case tcell.KeyRune:
		runes := []rune(ev.Str())
		for _, r := range runes {
			// If 'j' is pressed right after 'k', delete 'k' and switch to NORMAL mode
			if r == 'k' && m.lastChar == 'j' {
				e.DeleteChar()
				e.SetMode("NORMAL")
				return
			}
			e.InsertChar(r)
			m.lastChar = r
		}
	default:
		// Reset tracking fo non-rune keys (like arrows, tab, etc.)
		m.lastChar = 0
	}
}
