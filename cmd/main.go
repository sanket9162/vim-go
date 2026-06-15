package main

import (
	"log"

	"github.com/sanket9162/vim-go/internal/editor"
	"github.com/sanket9162/vim-go/internal/ui"
)

func main() {
	s, err := ui.NewScreen()
	if err != nil {
		log.Fatalf("%+V", err)
	}
	defer s.Close()

	e := editor.NewEditor(s)
	e.Run()
}
