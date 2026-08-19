package ui

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/Ashoke15/notes-app/internal/vault"
)

func (m Model) enter() (tea.Model, tea.Cmd) {

	if m.CurrentFile != nil {
		return m, nil
	}

	if m.State == stateConfirmDelete {
		return m, nil
	}

	if m.State == stateListing {
		item, ok := m.List.SelectedItem().(Item)
		if ok {
			f, content, err := vault.Open(vaultDir, item.title)
			if err != nil {
				m.StatusMsg = "Error: Cannot open file"
				return m, nil
			}

			m.NoteTextArea.SetValue(string(content))
			m.CurrentFile = f
			m.State = stateEditing
		}

		return m, nil
	}

	fileName := m.NewFileInput.Value()
	if fileName != "" {
		f, err := vault.Creat(vaultDir, fileName)
		if err != nil {
			if errors.Is(err, vault.ErrAlredyExists) {
				m.StatusMsg = "A note with that name alredy exits"
			} else {
				m.StatusMsg = " cannot crea file"
			}
			return m, nil
		}

		m.CurrentFile = f
		m.State = stateEditing
		m.NoteTextArea.SetValue("")
	}

	return m, nil
}

func (m Model) esc() (tea.Model, tea.Cmd) {

	if m.State == stateConfirmDelete {
		m.PendingDelete = ""
		m.State = stateListing

		return m, nil	
	}

	if m.State == stateEditing && m.CurrentFile != nil {
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

func (m Model) list() (tea.Model, tea.Cmd) {
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

	if err := vault.Save(m.CurrentFile, m.NoteTextArea.Value()); err != nil {
		m.StatusMsg = "Error: cannot save massege"
		return m, nil
	}

	m.CurrentFile = nil
	m.NoteTextArea.SetValue("")
	m.StatusMsg = "Not save successful"
	m.State = stateIdle

	return m, nil
}

func (m Model) startDelet() (tea.Model, tea.Cmd) {
	if m.State != stateListing || m.List.FilterState() == list.Filtering {
		return m, nil
	}

	item, ok := m.List.SelectedItem().(Item)
	if !ok {
		return m, nil
	}

	m.PendingDelete = item.Title()
	m.State = stateConfirmDelete

	return m, nil	
}

func (m Model) confirmDelet() (tea.Model, tea.Cmd) {
	if err := vault.Delet(vaultDir, m.PendingDelete); err != nil {
		m.StatusMsg = "Error: Cannot delet note"
	} else {
		m.StatusMsg = fmt.Sprintf("Deleted: %s",m.PendingDelete)
	}

	items, err := listFile()
	if err != nil {
		items = []list.Item{}
	}
	m.List.SetItems(items)

	m.PendingDelete = ""
	m.State = stateListing

	return m, nil	
}

func (m Model) cancelDelet() (tea.Model, tea.Cmd) {
	m.PendingDelete = ""
	m.State = stateListing

	return m, nil
}

func listFile() ([]list.Item, error) {
	notes, err := vault.List(vaultDir)
	if err != nil {
		return nil, err
	}

	items := make([]list.Item, 0, len(notes))
	for _, n := range notes {
		items = append(items, Item{
			title:       n.Name,
			description: fmt.Sprintf("Modified: %s", n.ModTime.Format("2006-01-02 15:04")),
		})
	}

	return items, nil
}
