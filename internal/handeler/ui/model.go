package ui

import (
	"os"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

type viewstate int

const (
	stateIdle viewstate = iota
	stateNewFile
	stateEditing
	stateListing
	stateRename
	stateConfirmDelete
)

type Model struct {
	NewFileInput  textinput.Model
	RenameInput   textinput.Model
	CurrentFile   *os.File
	NoteTextArea  textarea.Model
	Preview       viewport.Model
	PreviewOn     bool
	Help          help.Model
	Theme         string
	Style         Styles
	DarkBG        bool
	Width, Height int
	List          list.Model
	StatusMsg     string
	State         viewstate
	PendingDelete string
	PendingRename string
	Dirty         bool
	autosavegen   int
}

type Styles struct {
	app           lipgloss.Style
	title         lipgloss.Style
	statusMessage lipgloss.Style
	danger        lipgloss.Style
	confirmBox    lipgloss.Style
	confirmHint   lipgloss.Style
	keyHint       lipgloss.Style
	renameBox     lipgloss.Style
	noteBox       lipgloss.Style
}

type Item struct {
	title       string
	description string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.description }
func (i Item) FilterValue() string { return i.title }

func (m *Model) UpdateListProperties() {
	// Update list size.
	h, _ := m.Style.app.GetFrameSize()
	m.List.SetSize(m.Width-h, m.Height-4)

	//update the preview size
	m.Preview.SetWidth(m.Width - h)
	m.Preview.SetHeight(m.Height - 4)

	m.Help.SetWidth(m.Width - h)

	m.NoteTextArea.SetWidth(m.Width - h - 4)
	m.NoteTextArea.SetHeight(m.Height - 8)

	// Update the model and list styles.
	m.Style = NewStyles(m.DarkBG)
	m.List.Styles.Title = m.Style.title
}
