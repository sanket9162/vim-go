package mode

import "github.com/gdamore/tcell/v3"

type InsertMode struct{}

func (m *InsertMode) Name() string { return "INSERT" }

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
