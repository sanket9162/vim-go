package main

import (
	"log"

	"github.com/gdamore/tcell/v3"
)

func main() {

	// initialize the screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+V", err)
	}
	if err = s.Init(); err != nil {
		log.Fatalf("%+V", err)
	}
	defer s.Fini()

}
