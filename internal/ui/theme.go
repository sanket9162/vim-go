package ui

import (
	"encoding/json"
	"os"

	"github.com/gdamore/tcell/v3"
)

// Theme struct the custom theme configuration loaded from a json file.
type Theme struct {
	Name   string            `json:"name"`
	Editor EditorColors      `json:"editor"`
	Syntax map[string]string `json:"syntax"`
}

// EditorColors represents the custom style specification for the editro TUI components.
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

// LoadedTheme holds the resolved tcell.Color and tcell.Sytle variants for quick rendering access.
type LoadedTheme struct {
	Name string

	//Editor UI Styles
	DefaultStyle       tcell.Style
	GutterStyle        tcell.Style
	SelectionStyle     tcell.Style
	SearchMatchStyle   tcell.Style
	SearchCurrentStyle tcell.Style
	StatusBarStyle     tcell.Style

	//Resolved Syntax Style
	SyntaxStyles map[string]tcell.Style
}

// LoadTheme reads and parses a JSON theme file.
func LoadTheme(filepath string) (*LoadedTheme, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var rawTheme Theme
	if err := json.Unmarshal(data, &rawTheme); err != nil {
		return nil, err
	}

	return NewLoadedTheme(rawTheme), nil

}

// NewLoadedTheme creates a LoadedTheme from a Theme.
func NewLoadedTheme(t Theme) *LoadedTheme {
	lt := &LoadedTheme{
		Name:         t.Name,
		SyntaxStyles: make(map[string]tcell.Style),
	}

	//Resolves all editor styles
	resolveStyle := func(fgStr, bgStr string) tcell.Style {
		style := tcell.StyleDefault
		if fgStr != "" {
			style = style.Foreground(tcell.GetColor(fgStr))
		}
		if bgStr != "" {
			style = style.Background(tcell.GetColor(bgStr))
		}
		return style
	}

	// Resolves all editor component styles
	lt.DefaultStyle = resolveStyle(t.Editor.Foreground, t.Editor.Background)
	lt.GutterStyle = resolveStyle(t.Editor.GutterForeground, t.Editor.GutterBackground)
	lt.SelectionStyle = resolveStyle(t.Editor.SelectionForeground, t.Editor.SelectionBackground)
	lt.SearchMatchStyle = resolveStyle(t.Editor.SearchMatchForeground, t.Editor.SearchMatchBackground)
	lt.SearchCurrentStyle = resolveStyle(t.Editor.SearchCurrentForeground, t.Editor.SearchCurrentBackground)
	lt.StatusBarStyle = resolveStyle(t.Editor.StatusBarForeground, t.Editor.StatusBarBackground)

	//Resolves all syntax token styles
	for tokenType, fgColor := range t.Syntax {
		lt.SyntaxStyles[tokenType] = resolveStyle(fgColor, t.Editor.Background)
	}

	return lt

}

// Global or Editor-specific helper to translate hex/names to tcell.Color
func ResolveColor(name string) tcell.Color {
	return tcell.GetColor(name)
}
