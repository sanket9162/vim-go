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
	MoveCursorToStartOfLine()
	MoveCursorToStartOfFile()
	MoveCursorToEndOfFile()
	MoveCursorToEndOfLine()
	MoveCursorToNextWord()
	InsertChar(r rune)
	InsertNewline()
	DeleteChar()
	DeleteUnderCursor()
	DeleteLine()
	DeleteWord()
	SetMode(name string)
	GetMode(name string) Mode
	QuitEditor()
	Paste(before bool)
	ExecuteCommand(cmd string)
	SaveFile()
}

// Mode represents a state of the editor with its own key handling logic.
type Mode interface {
	HandleKey(e EditorInterface, ev *tcell.EventKey)
	Name() string
}
