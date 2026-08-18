package ui

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory", err)
	}

	vaultDir = fmt.Sprintf("%s/.totion", homedir)
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Model) View() tea.View {

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).
		PaddingLeft(2).
		PaddingRight(2)

	welcome := style.Render("Welcome to Totion 🧠")

	help := "Ctrl+N: new file . Ctrl+L: list . Esc: back . Ctrl+S: save . Ctrl+Q: quit"

	view := ""
	if m.InputVisible {
		view = m.NewFileInput.View()
	}

	if m.CurrentFile != nil {
		view = m.NoteTextArea.View()
	}

	if m.ShowingList {
		view = m.List.View()
	}

	prit := fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)

	return tea.NewView(prit)
}
