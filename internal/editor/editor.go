package editor

import (
	"github.com/sanket9162/vim-go/internal/buffer"
	"github.com/sanket9162/vim-go/internal/mode"
	"github.com/sanket9162/vim-go/internal/ui"
)

type Editor struct {
	Buffer      *buffer.Buffer
	Cursor      *buffer.Cursor
	Screen      *ui.Screen
	CurrentMode mode.Mode
	modes       map[string]mode.Mode
	Quit        bool
}

func NewEditor(s *ui.Screen) *Editor {
	b := buffer.NewBuffer()
	e := &Editor{
		Buffer: b,
		Cursor: buffer.NewCursor(b),
		Screen: s,
		modes:  make(map[string]mode.Mode),
	}

	e.modes["NORMAL"] = &mode.NormalMode{}
	e.modes["INSERT"] = &mode.InsertMode{}
	e.CurrentMode = e.modes["NORMAL"]

	return e
}

func (e *Editor) SetMode(name string) {
	if m, ok := e.modes[name]; ok {
		e.CurrentMode = m
	}
}

func (e *Editor) MoveCursorLeft()  { e.Cursor.MoveLeft() }
func (e *Editor) MoveCursorRight() { e.Cursor.MoveRight() }
func (e *Editor) MoveCursorUp()    { e.Cursor.MoveUp() }
func (e *Editor) MoveCursorDown()  { e.Cursor.MoveDown() }

func (e *Editor) InsertChar(r rune) {
	e.Buffer.InsertChar(e.Cursor.Row(), e.Cursor.Col(), r)
	e.Cursor.SetPos(e.Cursor.Col()+1, e.Cursor.Row())
}

func (e *Editor) InsertNewline() {
	e.Buffer.InsertNewline(e.Cursor.Row(), e.Cursor.Col())
	e.Cursor.SetPos(0, e.Cursor.Row()+1)
}

func (e *Editor) DeleteChar() {
	row, col := e.Buffer.DeleteChar(e.Cursor.Row(), e.Cursor.Col())
	e.Cursor.SetPos(col, row)
}

func (e *Editor) QuitEditor() {
	e.Quit = true
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
	status := "-- " + e.CurrentMode.Name() + " --"
	e.Screen.DrawText(0, h-1, status)

	e.Screen.ShowCursor(e.Cursor.Col(), e.Cursor.Row())
	e.Screen.Show()
}
