package mode

import "github.com/gdamore/tcell/v3"

// SearchMode handles input when the user is typing a query after pressing '/'
type SearchMode struct {
	Query string
}

// New returns the search text prefixed with a slash
func (m *SearchMode) Name() string {
	return "/" + m.Query
}

func (m *SearchMode) HandleKey(e EditorInterface, ev *tcell.EventKey) {}
