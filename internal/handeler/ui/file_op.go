package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func (m Model) enter() (tea.Model, tea.Cmd) {

	if m.CurrentFile != nil {
		return m, nil
	}

	if m.State == stateListing {
		item, ok := m.List.SelectedItem().(Item)
		if ok {
			notePath := filepath.Join(vaultDir, item.Title())
			content, err := os.ReadFile(notePath)
			if err != nil {
				m.StatusMsg = "Error: Cannot read file"
				return m, nil
			}

			m.NoteTextArea.SetValue(string(content))

			f, err := os.OpenFile(notePath, os.O_RDWR, 0644)
			if err != nil {
				m.StatusMsg = "Error: Cannot open file for editing"
				return m, nil
			}

			m.CurrentFile = f
			m.State = stateEditing

		}

		return m, nil
	}

	fileName := m.NewFileInput.Value()
	if fileName != "" {
		notePath := filepath.Join(vaultDir, fileName+".md")
		if _, err := os.Stat(notePath); err == nil {
			m.StatusMsg = "A note with that name already exists"
			return m, nil
		}

		f, err := os.Create(notePath)
		if err != nil {
			m.StatusMsg = "Error: Cannot creat file"
			return m, nil
		}

		m.CurrentFile = f
		m.State = stateEditing
		m.NewFileInput.SetValue("")
	}
	return m, nil
}

func(m Model) esc() (tea.Model, tea.Cmd){
	if m.State == stateListing && m.List.FilterState() == list.Filtering {
		return m, nil
	}

	if m.State == stateEditing && m.CurrentFile != nil{
		m.CurrentFile.Close()
		m.CurrentFile = nil
		m.NoteTextArea.SetValue("")
	}

	if m.State == stateNewFile {
		m.NewFileInput.SetValue("")
	}

	m.State = stateIdle
	m.StatusMsg = ""

	return m, nil
}

func (m Model) list() (tea.Model, tea.Cmd){
	if m.State != stateIdle {
		return m, nil
	}

	notelist, err := listFile()
	if err != nil {
		notelist = []list.Item{}
	}
	m.List.SetItems(notelist)
	m.State = stateListing

	return m, nil
}

func (m Model) save() (tea.Model, tea.Cmd) {
	if m.CurrentFile == nil {
		return m, nil
	}

	if err := m.CurrentFile.Truncate(0); err != nil {
		m.StatusMsg = "can not save file :("
		return m, nil
	}

	if _, err := m.CurrentFile.Seek(0, 0); err != nil {
		m.StatusMsg = "can not seek the file :("
		return m, nil
	}

	if _, err := m.CurrentFile.WriteString(m.NoteTextArea.Value()); err != nil {
		m.StatusMsg = "can not write the file :("
		return m, nil
	}

	if err := m.CurrentFile.Close(); err != nil {
		m.StatusMsg = "can not close the file."
		return m, nil
	}

	m.CurrentFile = nil
	m.NoteTextArea.SetValue("")
	m.StatusMsg = "Note Save successfully"
	m.State = stateIdle

	return m, nil
}


func listFile() ([]list.Item, error) {
	items := make([]list.Item, 0)

	entries, err := os.ReadDir(vaultDir)
	if err != nil {
		return nil, fmt.Errorf("error reading note list %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			modTime := info.ModTime().Format("2006-01-02 15:04")

			items = append(items, Item{
				title:       entry.Name(),
				description: fmt.Sprintf("modified: %s", modTime),
			})

		}
	}

	return items, nil
}
