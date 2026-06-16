package mode

import (
	"github.com/gdamore/tcell/v3"
)

// CommandMode handles the state where the user is typing a command after pressing ':'.
type CommandMode struct {
	Command string
}

// Name returns the current command string prefixed with a colon.
func (m *CommandMode) Name() string {
	return ":" + m.Command
}

// HandleKey handles key events while in command mode.
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
		m.Command += ev.Str()
	}
}
