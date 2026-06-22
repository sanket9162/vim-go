package mode

import (
	"github.com/gdamore/tcell/v3"
)

// SearchMode handles input when the user is typing a query after pressing '/'.
type SearchMode struct {
	Query string
}

// Name returns the search text prefixed with a slash.
func (m *SearchMode) Name() string {
	return "/" + m.Query
}

// HandleKey processes keystrokes in search mode.
func (m *SearchMode) HandleKey(e EditorInterface, ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		m.Query = ""
		e.PerformSearch("") // Clear highlights
		e.SetMode("NORMAL")
	case tcell.KeyEnter:
		e.PerformSearch(m.Query)
		m.Query = ""
		e.SetMode("NORMAL")
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if len(m.Query) > 0 {
			m.Query = m.Query[:len(m.Query)-1]
		} else {
			e.SetMode("NORMAL")
		}
	case tcell.KeyRune:
		m.Query += ev.Str()
	}
}
