package editor

import (
	"fmt"
	"strings"

	"github.com/sanket9162/vim-go/internal/buffer"
	"github.com/sanket9162/vim-go/internal/highlight"
	"github.com/sanket9162/vim-go/internal/mode"
	"github.com/sanket9162/vim-go/internal/ui"
)

// Selection represents a 2D text selection range.
type Selection struct {
	StartCol, StartRow int
	EndCol, EndRow     int
	Active             bool
}

// Editor is the main controller that ties the buffer, cursor, and screen together.
type Editor struct {
	Buffer        *buffer.Buffer
	Cursor        *buffer.Cursor
	Screen        *ui.Screen
	Viewport      *ui.Viewport
	CurrentMode   mode.Mode
	modes         map[string]mode.Mode
	Quit          bool
	FileName      string
	Selection     Selection
	Clipboard     string
	SearchQuery   string
	SearchResults []SearchMatch
	SearchIndex   int
	Theme         *ui.LoadedTheme
}

// SearchMatch represents a 2D text coordinate range for a search result.
type SearchMatch struct {
	Row int
	Col int
	Len int
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
		Theme:    ui.DefaultTheme(),
	}

	e.modes["NORMAL"] = &mode.NormalMode{}
	e.modes["INSERT"] = &mode.InsertMode{}
	e.modes["COMMAND"] = &mode.CommandMode{}
	e.modes["VISUAL"] = &mode.VisualMode{}
	e.modes["SEARCH"] = &mode.SearchMode{}
	e.CurrentMode = e.modes["NORMAL"]

	return e
}

// SetMode changes the editor's current input mode.
func (e *Editor) SetMode(name string) {
	if m, ok := e.modes[name]; ok {
		// save snapshot before entering INSERT mode.
		if name == "INSERT" && e.CurrentMode.Name() != "INSERT" {
			e.Buffer.SaveSnapshot(e.Cursor.Col(), e.Cursor.Row())
		}
		e.CurrentMode = m
		if name == "VISUAL" {
			if vm, ok := m.(*mode.VisualMode); ok {
				vm.InitSelection(e.Cursor.Col(), e.Cursor.Row())
				e.UpdateSelection(e.Cursor.Col(), e.Cursor.Row(), e.Cursor.Col(), e.Cursor.Row())
			}
		} else if name == "NORMAL" {
			e.ClearSelection()
			row := e.Cursor.Row()
			col := e.Cursor.Col()
			lineLen := e.Buffer.LineLength(row)
			if col >= lineLen && lineLen > 0 {
				e.Cursor.SetPos(lineLen-1, row)
			}
		}
	}
}

// GetMode returns the mode with the given name.
func (e *Editor) GetMode(name string) mode.Mode {
	return e.modes[name]
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

func (e *Editor) MoveCursorToEndOfLine() {
	e.Cursor.MoveToEndOfLine()
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
		tokenTypes := highlight.TokenizeLine(string(line), highlight.GoRules)
		for x := 0; x < e.Viewport.Width; x++ {
			bufferCol := x + e.Viewport.OffsetX
			if bufferCol >= len(line) {
				break
			}
			//Apply theme style for the token type
			style := e.Theme.GetTokenStyle(tokenTypes[bufferCol])
			if e.isSelected(bufferCol, bufferRow) {
				style = e.Theme.SelectionStyle
			} else if isMatch, isCurrent := e.isSearchMatch(bufferCol, bufferRow); isMatch {
				if isCurrent {
					style = e.Theme.SearchCurrentStyle
				} else {
					style = e.Theme.SearchMatchStyle
				}
			}
			e.Screen.SetContent(x+gutterWidth, y, line[bufferCol], nil, style)
		}
	}

	// Draw the enhanced status bar
	e.renderStatusBar(gutterWidth)

	visualX := e.Cursor.Col() - e.Viewport.OffsetX + gutterWidth
	visualY := e.Cursor.Row() - e.Viewport.OffsetY

	if strings.HasPrefix(e.CurrentMode.Name(), ":") || strings.HasPrefix(e.CurrentMode.Name(), "/") {
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
	isCmdOrSearch := strings.HasPrefix(modeName, ":") || strings.HasPrefix(modeName, "/")

	if isCmdOrSearch {
		// Just draw the query/command at the bottom row (overwrite everything on that row)
		e.Screen.DrawText(0, h-1, modeName+strings.Repeat(" ", w-len(modeName)))
		return
	}

	modeName = "--" + modeName + "--"
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

// GetCursorPos returns the current 2D cursor position.
func (e *Editor) GetCursorPos() (col, row int) {
	return e.Cursor.Col(), e.Cursor.Row()
}

// UpdateSelection sets the selection range parameters.
func (e *Editor) UpdateSelection(startCol, startRow, endCol, endRow int) {
	e.Selection = Selection{
		StartCol: startCol,
		StartRow: startRow,
		EndCol:   endCol,
		EndRow:   endRow,
		Active:   true,
	}
}

// ClearSelection deactivates text selection.
func (e *Editor) ClearSelection() {
	e.Selection.Active = false
}

// YankSelection copies the selected text into the editor's clipboard.
func (e *Editor) YankSelection() {
	if !e.Selection.Active {
		return
	}
	e.Clipboard = e.getSelectedText()
}

// DeleteSelection cuts the selected text from the buffer and positions the cursor at selection start.
func (e *Editor) DeleteSelection() {
	if !e.Selection.Active {
		return
	}
	e.Buffer.SaveSnapshot(e.Cursor.Col(), e.Cursor.Row())
	rStart, cStart, rEnd, cEnd := e.getNormalizedSelection()
	newRow, newCol := e.Buffer.DeleteRange(rStart, cStart, rEnd, cEnd)
	e.Cursor.SetPos(newCol, newRow)
}

func (e *Editor) getNormalizedSelection() (rStart, cStart, rEnd, cEnd int) {
	rStart, rEnd = e.Selection.StartRow, e.Selection.EndRow
	cStart, cEnd = e.Selection.StartCol, e.Selection.EndCol
	if rStart > rEnd || (rStart == rEnd && cStart > cEnd) {
		rStart, rEnd = rEnd, rStart
		cStart, cEnd = cEnd, cStart
	}
	return
}

func (e *Editor) getSelectedText() string {
	if !e.Selection.Active {
		return ""
	}
	rStart, cStart, rEnd, cEnd := e.getNormalizedSelection()
	var sb strings.Builder

	if rStart == rEnd {
		line := e.Buffer.GetLine(rStart)
		if cStart < len(line) {
			end := cEnd + 1
			if end > len(line) {
				end = len(line)
			}
			sb.WriteString(string(line[cStart:end]))
		}
		return sb.String()
	}

	// First line
	line := e.Buffer.GetLine(rStart)
	if cStart < len(line) {
		sb.WriteString(string(line[cStart:]))
	}
	sb.WriteRune('\n')

	// Middle lines
	for r := rStart + 1; r < rEnd; r++ {
		line = e.Buffer.GetLine(r)
		sb.WriteString(string(line))
		sb.WriteRune('\n')
	}

	// Last line
	line = e.Buffer.GetLine(rEnd)
	end := cEnd + 1
	if end > len(line) {
		end = len(line)
	}
	sb.WriteString(string(line[:end]))

	return sb.String()
}

func (e *Editor) isSelected(col, row int) bool {
	if !e.Selection.Active {
		return false
	}
	rStart, cStart, rEnd, cEnd := e.getNormalizedSelection()

	if row < rStart || row > rEnd {
		return false
	}
	if row > rStart && row < rEnd {
		return true
	}
	if rStart == rEnd {
		return col >= cStart && col <= cEnd
	}
	if row == rStart {
		return col >= cStart
	}
	if row == rEnd {
		return col <= cEnd
	}
	return false
}

// DeleteUnderCursor deletes the character directly under the cursor.
func (e *Editor) DeleteUnderCursor() {
	row := e.Cursor.Row()
	col := e.Cursor.Col()
	lineLen := e.Buffer.LineLength(row)
	if lineLen == 0 {
		return
	}

	if col >= lineLen {
		col = lineLen - 1
	}

	e.Buffer.SaveSnapshot(e.Cursor.Col(), e.Cursor.Row())
	e.Buffer.DeleteRange(row, col, row, col)

	// Shift cursor left if it was snapped past the new end of line
	newLineLen := e.Buffer.LineLength(row)
	if e.Cursor.Col() >= newLineLen && newLineLen > 0 {
		e.Cursor.SetPos(newLineLen-1, row)
	}
}

// DeleteLine deletes the entire current line including its trailing newline.
func (e *Editor) DeleteLine() {
	row := e.Cursor.Row()
	lineLen := e.Buffer.LineLength(row)

	e.Buffer.SaveSnapshot(e.Cursor.Col(), e.Cursor.Row())
	e.Buffer.DeleteRange(row, 0, row, lineLen)

	// Adjust cursor to keep it in a valid row bounds
	totalLines := e.Buffer.LineCount()
	if row >= totalLines {
		row = totalLines - 1
	}
	if row < 0 {
		row = 0
	}
	e.Cursor.SetPos(0, row)
}

// DeleteWord deletes from the cursor position to the beginning of the next word.
func (e *Editor) DeleteWord() {
	row := e.Cursor.Row()
	col := e.Cursor.Col()
	line := e.Buffer.GetLine(row)
	if len(line) == 0 {
		// On an empty line, dw deletes the newline character
		e.Buffer.SaveSnapshot(e.Cursor.Col(), e.Cursor.Row())
		e.Buffer.DeleteRange(row, col, row, col)
		return
	}

	nextCol := col
	// Move past current word chars
	for nextCol < len(line) && isWordChar(line[nextCol]) {
		nextCol++
	}
	// Move past trailing spaces
	for nextCol < len(line) && !isWordChar(line[nextCol]) {
		nextCol++
	}

	e.Buffer.SaveSnapshot(e.Cursor.Col(), e.Cursor.Row())
	if nextCol == col {
		// If we're at the end of a line, delete the newline
		e.Buffer.DeleteRange(row, col, row, col)
	} else {
		e.Buffer.DeleteRange(row, col, row, nextCol-1)
	}
}

// Helper to check for word characters (matching cursor logic)
func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.' || r == '!' || r == ':' || r == '+'
}

func (e *Editor) Paste(before bool) {
	if e.Clipboard == "" {
		return
	}

	e.Buffer.SaveSnapshot(e.Cursor.Col(), e.Cursor.Row())
	row := e.Cursor.Row()
	col := e.Cursor.Col()
	isLine := strings.HasSuffix(e.Clipboard, "\n")

	if isLine {
		lines := strings.Split(e.Clipboard, "\n")
		// Remove empty trailing string from split
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}

		targetRow := row
		if !before {
			targetRow = row + 1
		}

		// Insert lines in order
		for i, lineText := range lines {
			e.Buffer.InsertLine(targetRow+i, lineText)
		}

		e.Cursor.SetPos(0, targetRow)

	} else {
		targetCol := col
		if !before {
			lineLen := e.Buffer.LineLength(row)
			if lineLen > 0 {
				targetCol = col + 1
			}
		}

		// Insert character-wise contents
		for i, r := range e.Clipboard {
			e.Buffer.InsertChar(row, targetCol+i, r)
		}

		e.Cursor.SetPos(targetCol+len(e.Clipboard)-1, row)

	}

}

func (e *Editor) Undo() {
	cCol, cRow := e.Cursor.Col(), e.Cursor.Row()
	if newCol, newRow, ok := e.Buffer.Undo(cCol, cRow); ok {
		e.Cursor.SetPos(newCol, newRow)
	}
}

func (e *Editor) Redo() {
	cCol, cRow := e.Cursor.Col(), e.Cursor.Row()
	if newCol, newRow, ok := e.Buffer.Redo(cCol, cRow); ok {
		e.Cursor.SetPos(newCol, newRow)
	}
}

// PerformSearch scans the buffer for string occurrences and populates highlights.
func (e *Editor) PerformSearch(query string) {
	e.SearchQuery = query
	e.SearchResults = []SearchMatch{}
	e.SearchIndex = -1

	if query == "" {
		return
	}

	for r := 0; r < e.Buffer.LineCount(); r++ {
		line := string(e.Buffer.GetLine(r))
		lowerLine := strings.ToLower(line)
		lowerQuery := strings.ToLower(query)

		start := 0
		for {
			idx := strings.Index(lowerLine[start:], lowerQuery)
			if idx == -1 {
				break
			}
			matchCol := start + idx
			e.SearchResults = append(e.SearchResults, SearchMatch{
				Row: r,
				Col: matchCol,
				Len: len(query),
			})
			start = matchCol + len(query)
			if len(query) == 0 {
				break
			}
		}
	}

	if len(e.SearchResults) > 0 {
		// Jump to first match at or after cursor position
		cursorRow := e.Cursor.Row()
		cursorCol := e.Cursor.Col()
		e.SearchIndex = 0
		for i, match := range e.SearchResults {
			if match.Row > cursorRow || (match.Row == cursorRow && match.Col >= cursorCol) {
				e.SearchIndex = i
				break
			}
		}
		e.JumpToMatch(e.SearchIndex)
	}
}

// JumpToMatch focuses the viewport/cursor on the selected match.
func (e *Editor) JumpToMatch(index int) {
	if index < 0 || index >= len(e.SearchResults) {
		return
	}
	match := e.SearchResults[index]
	e.Cursor.SetPos(match.Col, match.Row)
	e.Viewport.ScrollTo(match.Col, match.Row)
}

// SearchNext advances to the next search result.
func (e *Editor) SearchNext() {
	if len(e.SearchResults) == 0 {
		return
	}
	e.SearchIndex = (e.SearchIndex + 1) % len(e.SearchResults)
	e.JumpToMatch(e.SearchIndex)
}

// SearchPrev wraps back to the previous search result.
func (e *Editor) SearchPrev() {
	if len(e.SearchResults) == 0 {
		return
	}
	e.SearchIndex = (e.SearchIndex - 1 + len(e.SearchResults)) % len(e.SearchResults)
	e.JumpToMatch(e.SearchIndex)
}

// isSearchMatch determines if a character position is part of a match.
func (e *Editor) isSearchMatch(col, row int) (bool, bool) {
	if len(e.SearchResults) == 0 {
		return false, false
	}
	for i, match := range e.SearchResults {
		if match.Row == row && col >= match.Col && col < match.Col+match.Len {
			return true, i == e.SearchIndex
		}
	}
	return false, false
}
