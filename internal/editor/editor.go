package editor

import (
	"fmt"
	"strings"

	"github.com/sanket9162/vim-go/internal/buffer"
	"github.com/sanket9162/vim-go/internal/mode"
	"github.com/sanket9162/vim-go/internal/ui"
)

// Editor is the main controller that ties the buffer, cursor, and screen together.
type Editor struct {
	Buffer      *buffer.Buffer
	Cursor      *buffer.Cursor
	Screen      *ui.Screen
	Viewport    *ui.Viewport
	CurrentMode mode.Mode
	modes       map[string]mode.Mode
	Quit        bool
	FileName    string
}

// NewEditor initializes a new Editor instance.
func NewEditor(s *ui.Screen, filename string) *Editor {
	w, h := s.Size()
	b := buffer.NewBuffer()

	// Load file if provided
	if filename != "" {
		_ = b.Load(filename)
	}

	e := &Editor{
		Buffer:   b,
		Cursor:   buffer.NewCursor(b),
		Screen:   s,
		Viewport: ui.NewViewport(w, h),
		FileName: filename,
		modes:    make(map[string]mode.Mode),
	}

	e.modes["NORMAL"] = &mode.NormalMode{}
	e.modes["INSERT"] = &mode.InsertMode{}
	e.modes["COMMAND"] = &mode.CommandMode{}
	e.CurrentMode = e.modes["NORMAL"]

	return e
}

// SetMode changes the editor's current input mode.
func (e *Editor) SetMode(name string) {
	if m, ok := e.modes[name]; ok {
		e.CurrentMode = m
	}
}

// MoveCursorLeft moves the cursor one position to the left.
func (e *Editor) MoveCursorLeft() { e.Cursor.MoveLeft() }

// MoveCursorRight moves the cursor one position to the right.
func (e *Editor) MoveCursorRight() { e.Cursor.MoveRight() }

// MoveCursorUp moves the cursor one position up.
func (e *Editor) MoveCursorUp() { e.Cursor.MoveUp() }

// MoveCursorDown moves the cursor one position down.
func (e *Editor) MoveCursorDown() { e.Cursor.MoveDown() }

// InsertChar inserts a single character into the buffer at the cursor position.
func (e *Editor) InsertChar(r rune) {
	e.Buffer.InsertChar(e.Cursor.Row(), e.Cursor.Col(), r)
	e.Cursor.SetPos(e.Cursor.Col()+1, e.Cursor.Row())
}

// InsertNewline inserts a newline at the cursor position.
func (e *Editor) InsertNewline() {
	e.Buffer.InsertNewline(e.Cursor.Row(), e.Cursor.Col())
	e.Cursor.SetPos(0, e.Cursor.Row()+1)
}

// DeleteChar deletes the character before the cursor position.
func (e *Editor) DeleteChar() {
	row, col := e.Buffer.DeleteChar(e.Cursor.Row(), e.Cursor.Col())
	e.Cursor.SetPos(col, row)
}

// QuitEditor sets the flag to exit the editor.
func (e *Editor) QuitEditor() {
	e.Quit = true
}

func (e *Editor) MoveCursorToStartOfLine() {
	e.Cursor.MoveToStartOfLine()
}

func (e *Editor) MoveCursorToStartOfFile() {
	e.Cursor.MoveToStartOfFile()
}

func (e *Editor) MoveCursorToEndOfFile() {
	e.Cursor.MoveToEndOfFile()
}

func (e *Editor) MoveCursorToNextWord() {
	e.Cursor.MoveToNextWord()
}

// Render updates the visual state of the editor on the screen.
func (e *Editor) Render() {
	e.Screen.Clear()

	// Calculate line number gutter width (e.g., " 1" is 4 chars)
	totalLines := e.Buffer.LineCount()
	gutterWidth := len(fmt.Sprintf("%d", totalLines)) + 2

	// Adjust viewport width to account for the gutter
	screenWidth, screenHeight := e.Screen.Size()
	e.Viewport.Width = screenWidth - gutterWidth
	e.Viewport.Height = screenHeight - 1

	e.Viewport.ScrollTo(e.Cursor.Col(), e.Cursor.Row())

	for y := 0; y < e.Viewport.Height; y++ {
		bufferRow := y + e.Viewport.OffsetY
		if bufferRow >= totalLines {
			break
		}

		// Draw Line Number
		lineNumstr := fmt.Sprintf("%*d", gutterWidth-1, bufferRow+1)
		// Optional : Draw with a different style/color
		e.Screen.DrawText(0, y, lineNumstr)

		line := e.Buffer.GetLine(bufferRow)
		for x := 0; x < e.Viewport.Width; x++ {
			bufferCol := x + e.Viewport.OffsetX
			if bufferCol >= len(line) {
				break
			}
			e.Screen.DrawRune(x+gutterWidth, y, line[bufferCol])
		}
	}

	// Draw the enhanced status bar
	e.renderStatusBar(gutterWidth)

	visualX := e.Cursor.Col() - e.Viewport.OffsetX + gutterWidth
	visualY := e.Cursor.Row() - e.Viewport.OffsetY

	// if e.CurrentMode.Name()[0] == ':' {
	// 	visualX = len(e.CurrentMode.Name())
	// 	visualY = screenHeight - 1
	// }

	if strings.HasPrefix(e.CurrentMode.Name(), ":") {
		visualX = len(e.CurrentMode.Name())
		visualY = screenHeight - 1
	}

	e.Screen.ShowCursor(visualX, visualY)
	e.Screen.Show()
}

// SaveFile writes the current buffer content back to the associated file.
func (e *Editor) SaveFile() {
	if e.FileName != "" {
		_ = e.Buffer.Save(e.FileName)
	}
}

func (e *Editor) renderStatusBar(gutterWidth int) {
	w, h := e.Screen.Size()

	// Left side: Mode and Filename
	modeName := e.CurrentMode.Name()
	if !strings.HasPrefix(modeName, ":") {
		modeName = "--" + modeName + "--"
	}

	fileName := e.FileName
	if fileName == "" {
		fileName = "[No Name]"
	}

	leftStatus := fmt.Sprintf("%s %s", modeName, fileName)

	// Right side: Row and Column
	rightStatus := fmt.Sprintf("%d,%d", e.Cursor.Row()+1, e.Cursor.Col()+1)

	// Draw left part
	e.Screen.DrawText(0, h-1, leftStatus)

	// Draw right part (aligned to right)
	if w > len(leftStatus)+len(rightStatus)+2 {
		e.Screen.DrawText(w-len(rightStatus), h-1, rightStatus)
	}
}

func (e *Editor) ExecuteCommand(cmd string) {
	switch cmd {
	case "q":
		e.QuitEditor()
	case "w":
		e.SaveFile()
	case "wq":
		e.SaveFile()
		e.QuitEditor()

	}
}
