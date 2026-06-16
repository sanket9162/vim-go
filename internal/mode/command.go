package mode

import "github.com/gdamore/tcell/v3"

type CommandMode struct {
	Command string
}

func (m *CommandMode) Name() string { return ":" + m.Command }

func (m *CommandMode) HandleKey(e EditorInterface, ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		m.Command = ""
		e.SetMode("NORMAL")
	case tcell.KeyEnter:
		e.ExecuteCommand(m.Command)
		m.Command = ""
		e.SetMode("NORMAL")
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if len(m.Command) > 0 {
			m.Command = m.Command[:len(m.Command)-1]
		} else {
			e.SetMode("NORMAL")
		}
	case tcell.KeyRune:
		m.Command += string(ev.Rune())
	}
}
