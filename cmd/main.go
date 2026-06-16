package main

import (
	"log"
	"os"

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

	filename := ""
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	// Create and run the editor.
	e := editor.NewEditor(s, filename)
	e.Run()
}
