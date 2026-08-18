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
			return m, tea.Quit

		case "q":
			if !m.InputVisible && m.CurrentFile == nil && !m.ShowingList {
				return m, tea.Quit
			}

		case "esc":
			if m.InputVisible {
				m.InputVisible = false
			}

			if m.CurrentFile != nil {
				m.NoteTextArea.SetValue("")
				m.CurrentFile = nil
			}

			if m.ShowingList {
				if m.List.FilterState() == list.Filtering {
					break
				}

				m.ShowingList = false
			}

			return m, nil

		case "ctrl+l":
			//todo: showing list
			noteList, err := listFile()
			if err != nil{
				noteList = []list.Item{}
			}
			m.List.SetItems(noteList)
			m.ShowingList = true
			return m, nil

		case "ctrl+n":
			m.InputVisible = true
			return m, nil

		case "ctrl+s":
			return m.save()

		case "enter":
			return m.enter()
		}
	}

	if m.InputVisible {
		m.NewFileInput, cmd = m.NewFileInput.Update(msg)
	}

	if m.CurrentFile != nil {
		m.NoteTextArea, cmd = m.NoteTextArea.Update(msg)
	}

	if m.ShowingList {
		m.List, cmd = m.List.Update(msg)
	}

	return m, cmd
}
