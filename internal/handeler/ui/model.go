package ui

import (
	"os"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

type Model struct {
	NewFileInput  textinput.Model
	InputVisible  bool
	CurrentFile   *os.File
	NoteTextArea  textarea.Model
	Style         Styles
	DarkBG        bool
	Width, Height int
	List          list.Model
	ShowingList   bool
	StatusMsg     string
}

type Styles struct {
	app           lipgloss.Style
	title         lipgloss.Style
	statusMessage lipgloss.Style
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

	// Update the model and list styles.
	m.Style = NewStyles(m.DarkBG)
	m.List.Styles.Title = m.Style.title
}
