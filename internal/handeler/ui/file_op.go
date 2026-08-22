package ui

import (
	"errors"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/Ashoke15/notes-app/internal/vault"
	"github.com/charmbracelet/glamour"
)

func (m Model) enter() (tea.Model, tea.Cmd) {

	if m.CurrentFile != nil {
		return m, nil
	}

	if m.State == stateConfirmDelete {
		return m, nil
	}

	if m.State == stateRename {
		return m.summitRename()
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
			m.Dirty = false
			m.autosavegen = 0
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
		m.Dirty = false
		m.autosavegen = 0
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
		if m.Dirty {
			_ = vault.Save(m.CurrentFile, m.NoteTextArea.Value())
		} else {
			m.CurrentFile.Close()
		}
		m.CurrentFile = nil
		m.NoteTextArea.SetValue("")
	}

	if m.State == stateEditing && m.PreviewOn {
		m.PreviewOn = false
		return m, nil
	}

	if m.State == stateNewFile {
		m.NewFileInput.SetValue("")
	}

	if m.State == stateRename {
		return m.cancelRename()
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
	m.PreviewOn = false
	m.Dirty = false
	m.StatusMsg = "Note save successful"
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
		m.StatusMsg = fmt.Sprintf("Deleted: %s", m.PendingDelete)
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

func (m Model) startRename() (tea.Model, tea.Cmd) {
	if m.State != stateListing || m.List.FilterState() == list.Filtering {
		return m, nil
	}

	item, ok := m.List.SelectedItem().(Item)
	if !ok {
		return m, nil
	}

	m.PendingRename = item.Title()
	m.RenameInput.SetValue(strings.TrimSuffix(item.Title(), ".md"))
	m.RenameInput.Focus()
	m.State = stateRename

	return m, nil
}

func (m Model) summitRename() (tea.Model, tea.Cmd) {
	newName := strings.TrimSpace(m.RenameInput.Value())

	if newName == "" || newName == strings.TrimSuffix(m.PendingRename, ".md") {
		return m.cancelRename()
	}

	nameWithExt := newName + ".md"

	err := vault.Rename(vaultDir, m.PendingRename, nameWithExt)
	switch {
	case errors.Is(err, vault.ErrAlredyExists):
		m.StatusMsg = "A note with that name alredy exists"
		return m, nil
	case err != nil:
		m.StatusMsg = "Error: cannot rename note"
		return m, nil
	}

	m.StatusMsg = fmt.Sprintf("Rename to %s", nameWithExt)

	item, err := listFile()
	if err != nil {
		item = []list.Item{}
	}
	m.List.SetItems(item)

	m.PendingRename = ""
	m.RenameInput.SetValue("")
	m.State = stateListing

	return m, nil
}

func (m Model) cancelRename() (tea.Model, tea.Cmd) {
	m.PendingRename = ""
	m.RenameInput.SetValue("")
	m.State = stateListing

	return m, nil
}

func (m Model) togglePreview() (tea.Model, tea.Cmd) {
	if m.State != stateEditing {
		return m, nil
	}

	if m.PreviewOn {
		m.PreviewOn = false
		return m, nil
	}

	rendered, err := renderMarkDown(m.NoteTextArea.Value(), m.Preview.Width())
	if err != nil {
		m.StatusMsg = "Error: Cannot render preview"
		return m, nil
	}

	m.Preview.SetContent(rendered)
	m.Preview.GotoTop()
	m.PreviewOn = true

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

func renderMarkDown(content string, width int) (string, error) {
	if width < 20 {
		width = 20
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)

	if err != nil {
		return "", fmt.Errorf("Custom build render err = %w", err)
	}

	rendered, err := renderer.Render(content)
	if err != nil {
		return "", fmt.Errorf("Custom render markdown err = %w", err)
	}

	return rendered, nil
}
