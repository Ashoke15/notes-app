package ui

import (
	"fmt"
	"log"
	"os"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func (m Model) enter() (tea.Model, tea.Cmd) {

	if m.CurrentFile != nil {
		return m, nil
	}

	if m.ShowingList {
		item, ok := m.List.SelectedItem().(Item)
		if ok {
			filepath := fmt.Sprintf("%s/%s", vaultDir, item.Title())
			content, err := os.ReadFile(filepath)
			if err != nil {
				log.Printf("Error reading file: %v", err)
				return m, nil
			}

			m.NoteTextArea.SetValue(string(content))

			f, err := os.OpenFile(filepath, os.O_RDWR, 0644)
			if err != nil {
				log.Printf("error reading file: %v", err)
			}

			m.CurrentFile = f
			m.ShowingList = false

		}

		return m, nil
	}

	fileName := m.NewFileInput.Value()
	if fileName != "" {
		filepath := fmt.Sprintf("%s/%s.md", vaultDir, fileName)

		if _, err := os.Stat(filepath); err == nil {
			return m, nil
		}

		f, err := os.Create(filepath)
		if err != nil {
			return m, nil
		}

		m.CurrentFile = f
		m.InputVisible = false
		m.NewFileInput.SetValue("")
	}
	return m, nil
}

func listFile() ([]list.Item, error) {
	items := make([]list.Item, 0)

	entries, err := os.ReadDir(vaultDir)
	if err != nil {
		return nil, fmt.Errorf("error reading note list %w",err)
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

func (m *Model) save() (tea.Model, tea.Cmd) {
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

	return m, nil
}
