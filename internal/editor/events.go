package editor

import (
	"github.com/gdamore/tcell/v3"
)

// Run starts the main event loop of the editor.
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

// handleKey handles the key press events from the user.
func (e *Editor) handleKey(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyCtrlC {
		e.Quit = true
		return
	}

	e.CurrentMode.HandleKey(e, ev)
}
