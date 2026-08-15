package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	newFileInput textinput.Model
	inputVisible bool
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+n":
			m.inputVisible = true
			return m, nil
		}
	}

	if m.inputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() tea.View {

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).
		PaddingLeft(2).
		PaddingRight(2)

	welcome := style.Render("Welcome to Totion 🧠")

	help := "Ctrl+N: new file . Ctrl+L: list . Esc: back/save . Ctrl+S: save . Ctrl+Q: quit"

	view := ""
	if m.inputVisible {
		view = m.newFileInput.View()
	}

	prit := fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)

	return tea.NewView(prit)
}

func initializeModel() model {

	//initialise new file input
	ti := textinput.New()
	ti.Placeholder = "What would like to call it"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(50)

	return model{
		newFileInput: ti,
		inputVisible: false,
	}
}

func main() {

	p := tea.NewProgram(initializeModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
