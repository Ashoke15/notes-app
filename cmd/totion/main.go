package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/Ashoke15/notes-app/internal/handeler/ui"
)
func main() {

	p := tea.NewProgram(ui.InitializeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
