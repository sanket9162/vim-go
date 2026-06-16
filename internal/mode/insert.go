package mode

import "github.com/gdamore/tcell/v3"

type InsertMode struct{}

// Name returns the name of the insert mode.
func (m *InsertMode) Name() string { return "INSERT" }

// HandleKey handles the key press events in insert mode.
func (m *InsertMode) HandleKey(e EditorInterface, ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		e.SetMode("NORMAL")
	case tcell.KeyEnter:
		e.InsertNewline()
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		e.DeleteChar()
	case tcell.KeyRune:
		runes := []rune(ev.Str())
		for _, r := range runes {
			e.InsertChar(r)
		}
	}
}
