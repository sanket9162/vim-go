package mode

import "github.com/gdamore/tcell/v3"

// VisualEditorInterface extends EditorInterface with selection
type VisualEditorInterface interface {
	EditorInterface
	GetCursorPos() (col, row int)
	UpdateSelection(startCol, startRow, endCol, endRow int)
	ClearSelection()
	YankSelection()
	DeleteSelection()
}

// VisualMode tracks selection anchor point and handles key inputs.
type VisualMode struct {
	startCol int
	startRow int
}

func (m *VisualMode) Name() string { return "VISUAL" }

// InitSelection sets the starting coordinate when entering visual mode.
func (m *VisualMode) InitSelection(col, row int) {
	m.startCol = col
	m.startRow = row
}

func (m *VisualMode) HandleKey(e EditorInterface, ev *tcell.EventKey) {
	// Verify the editor supperts visual selection operation
	vEdit, ok := e.(VisualEditorInterface)
	if !ok {
		if ev.Key() == tcell.KeyEscape {
			e.SetMode("NORMAL")
		}
		return
	}

	switch ev.Key() {
	case tcell.KeyEscape:
		vEdit.ClearSelection()
		vEdit.SetMode("NORMAL")
		return

	}

	if ev.Key() == tcell.KeyRune {
		switch ev.Str() {
		// Navigation keys
		case "h":
			e.MoveCursorLeft()
		case "l":
			e.MoveCursorRight()
		case "j":
			e.MoveCursorDown()
		case "k":
			e.MoveCursorUp()
		case "0":
			e.MoveCursorToStartOfLine()
		case "$":
			e.MoveCursorToEndOfLine()
		case "w":
			e.MoveCursorToNextWord()
			// Operations on selection
		case "y": // Yank (copy) selection
			vEdit.YankSelection()
			vEdit.ClearSelection()
			e.SetMode("NORMAL")
			return
		case "d", "x": // Delete (cut) selection
			vEdit.DeleteSelection()
			vEdit.ClearSelection()
			e.SetMode("NORMAL")
			return
		}
		// Recalculate selection area after movement
		curCol, curRow := vEdit.GetCursorPos()
		vEdit.UpdateSelection(m.startCol, m.startRow, curCol, curRow)

	}
}
