package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

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

	help := "Ctrl+N: new file . Ctrl+L: list . d: delete . Esc/q: back . Ctrl+S: save "

	view := ""

	switch m.State {
	case stateNewFile:
		view = m.NewFileInput.View()

	case stateEditing:
		view = m.NoteTextArea.View()

	case stateListing:
		view = m.List.View()

	case stateIdle:
		view = "Press Ctrl+N for new note, or Ctrl+L to view all notes."

	case stateConfirmDelete:
		warning := m.Style.danger.Bold(true).Render("⚠ Delete note?")
		name := m.Style.danger.Render(fmt.Sprintf("\"%s\"", m.PendingDelete))
		note := m.Style.confirmHint.Render("This cannot be undone.")

		yKey := m.Style.keyHint.Render(" y ")
		nKey := m.Style.keyHint.Render(" n ")
		actions := fmt.Sprintf("%s confirm    %s cancel", yKey, nKey)

		body := fmt.Sprintf("%s\n%s\n\n%s\n\n%s", warning, name, note, actions)
		view = m.Style.confirmBox.Render(body)
	}

	prit := fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)

	if m.StatusMsg != "" {
		renderdStatus := m.Style.statusMessage.Render(m.StatusMsg)
		prit = fmt.Sprintf("%s\n\n%s", prit, renderdStatus)
	}

	return tea.NewView(prit)
}
