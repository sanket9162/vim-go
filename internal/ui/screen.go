package ui

import "github.com/gdamore/tcell/v3"

type Screen struct {
	tcell.Screen
	defStyle tcell.Style
}

func NewScreen() (*Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := s.Init(); err != nil {
		return nil, err
	}

	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(style)
	return &Screen{
		Screen:   s,
		defStyle: style,
	}, nil
}

// DrawText draws a string of text horizontally on the terminal screen.
func (s *Screen) DrawText(x, y int, text string) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, s.defStyle)
	}
}

// DrawRune draws a single rune at specific coordinates on the terminal screen.
func (s *Screen) DrawRune(x, y int, r rune) {
	s.SetContent(x, y, r, nil, s.defStyle)
}

// Close clears the terminal screen before exiting.
func (s *Screen) Close() {
	s.Fini()
}
