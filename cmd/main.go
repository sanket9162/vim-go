package main

import (
	"log"

	"github.com/sanket9162/vim-go/internal/editor"
	"github.com/sanket9162/vim-go/internal/ui"
)

func main() {
	// Initialize the TUI screen.
	s, err := ui.NewScreen()
	if err != nil {
		log.Fatalf("%+V", err)
	}
	defer s.Close()

	// Create and run the editor.
	e := editor.NewEditor(s)
	e.Run()
}
