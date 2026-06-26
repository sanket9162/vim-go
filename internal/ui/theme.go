package ui

import "github.com/gdamore/tcell/v3"

type Theme struct {
	Name   string            `json:"name"`
	Editor EditorColors      `json:"editor"`
	Syntax map[string]string `json:"syntax"`
}

type EditorColors struct {
	Background              string `json:"background"`
	Foreground              string `json:"foreground"`
	GutterBackground        string `json:"gutter_background"`
	GutterForeground        string `json:"gutter_foreground"`
	SelectionBackground     string `json:"selection_background"`
	SelectionForeground     string `json:"selection_foreground"`
	SearchMatchBackground   string `json:"search_match_background"`
	SearchMatchForeground   string `json:"search_match_foreground"`
	SearchCurrentBackground string `json:"search_current_background"`
	SearchCurrentForeground string `json:"search_current_foreground"`
	StatusBarBackground     string `json:"status_bar_background"`
	StatusBarForeground     string `json:"status_bar_foreground"`
}

// Global or Editor-specific helper to translate hex/names to tcell.Color
func ResolveColor(name string) tcell.Color {
	return tcell.GetColor(name)
}
