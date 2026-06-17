package mode

import (
	"github.com/gdamore/tcell/v3"
)

// EditorInterface defines the operations that modes can perform on the editor.
type EditorInterface interface {
	MoveCursorLeft()
	MoveCursorRight()
	MoveCursorUp()
	MoveCursorDown()
	MoveCursorToStartOfFile()
	InsertChar(r rune)
	InsertNewline()
	DeleteChar()
	SetMode(name string)
	QuitEditor()
	ExecuteCommand(cmd string)
	SaveFile()
}

// Mode represents a state of the editor with its own key handling logic.
type Mode interface {
	HandleKey(e EditorInterface, ev *tcell.EventKey)
	Name() string
}
