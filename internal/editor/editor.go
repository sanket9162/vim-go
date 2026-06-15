package editor

import (
	"github.com/sanket9162/vim-go/internal/buffer"
	"github.com/sanket9162/vim-go/internal/ui"
)

type Editor struct {
	Buffer *buffer.Buffer
	Cursor *buffer.Cursor
	Screen *ui.Screen
	Mode   string
	Quit   bool
}

func NewEditor(s *ui.Screen) *Editor {
	b := buffer.NewBuffer()
	return &Editor{
		Buffer: b,
		Cursor: buffer.NewCursor(b),
		Screen: s,
		Mode:   "NORMAL",
	}
}

// Render tells the UI to draw the buffer and status line
func (e *Editor) Render() {
	e.Screen.Clear()

	// Draw Buffer
	for y, line := range e.Buffer.Lines {
		for x, r := range line {
			e.Screen.DrawRune(x, y, r)
		}
	}

	// Draw Status Line
	_, h := e.Screen.Size()
	status := "-- " + e.Mode + " --"
	e.Screen.DrawText(0, h-1, status)

	e.Screen.ShowCursor(e.Cursor.Col(), e.Cursor.Row())
	e.Screen.Show()
}
