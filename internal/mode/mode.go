package mode

import (
	"github.com/gdamore/tcell/v3"
)

type EditorInterface interface {
	MoveCursorLeft()
	MoveCursorRight()
	MoveCursorUp()
	MoveCursorDown()
	InsertChar(r rune)
	InsertNewline()
	DeleteChar()
	SetMode(name string)
	QuitEditor()
}

type Mode interface {
	HandleKey(e EditorInterface, ev *tcell.EventKey)
	Name() string
}
