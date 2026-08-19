package ui

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.DarkBG = msg.IsDark()
		m.UpdateListProperties()
		return m, nil

	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height
		m.UpdateListProperties()
		return m, nil

	case tea.KeyPressMsg:
		m.StatusMsg = ""

		switch msg.String() {
		case "ctrl+c":
			if m.CurrentFile != nil {
				m.CurrentFile.Close()
			}
			return m, tea.Quit

		case "q":
			if m.State == stateIdle && m.CurrentFile == nil {
				return m, tea.Quit
			}

		case "d":
			if m.State == stateListing && m.List.FilterState() != list.Filtering {
				return m.startDelet()
			}

		case "y":
			if m.State == stateConfirmDelete {
				return m.confirmDelet()
			}

		case "n":
			if m.State == stateConfirmDelete {
				return m.cancelDelet()
			}

		case "esc":
			if m.State == stateListing && m.List.FilterState() == list.Filtering {
				break
			}

			return m.esc()

		case "ctrl+l":
			return m.list()

		case "ctrl+n":
			if m.State != stateIdle {
				return m, nil
			}
			m.State = stateNewFile
			return m, nil

		case "ctrl+s":
			return m.save()

		case "enter":
			return m.enter()
		}
	}

	if m.State == stateNewFile {
		m.NewFileInput, cmd = m.NewFileInput.Update(msg)
	}

	if m.CurrentFile != nil {
		m.NoteTextArea, cmd = m.NoteTextArea.Update(msg)
	}

	if m.State == stateListing {
		m.List, cmd = m.List.Update(msg)
	}

	return m, cmd
}
