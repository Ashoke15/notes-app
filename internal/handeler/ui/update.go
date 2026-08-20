package ui

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		if m.Theme == "auto" {
			m.DarkBG = msg.IsDark()
		}
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

		case "r":
			if m.State == stateListing && m.List.FilterState() != list.Filtering {
				return m.startRename()
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

		case "ctrl+p":
			return m.togglePreview()

		case "enter":
			return m.enter()

		case "?":
			switch m.State {
			case stateNewFile, stateRename:
				
			case stateEditing:
				if !m.PreviewOn{
					break
				}
				fallthrough

			default:
				m.Help.ShowAll = !m.Help.ShowAll
				return m, nil	
			}
		}
	}

	if m.State == stateNewFile {
		m.NewFileInput, cmd = m.NewFileInput.Update(msg)
	}

	if m.CurrentFile != nil {
		if m.PreviewOn {
			m.Preview, cmd = m.Preview.Update(msg)
		} else {
			m.NoteTextArea, cmd = m.NoteTextArea.Update(msg)
		}
	}

	if m.State == stateListing {
		m.List, cmd = m.List.Update(msg)
	}

	if m.State == stateRename {
		m.RenameInput, cmd = m.RenameInput.Update(msg)
	}

	return m, cmd
}
