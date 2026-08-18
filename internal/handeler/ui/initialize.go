package ui

import (
	"log"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

var (
	cursorColor = lipgloss.Color("205")

	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))

	focusedPlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))

	focusedPromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))

	vaultDir string
)

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory", err)
	}

	vaultDir = filepath.Join(homedir, ".totion")
}

func InitializeModel() Model {

	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatal(err)
	}

	//initialise new file input
	ti := textinput.New()
	ti.Placeholder = "What would like to call it"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(50)

	s := ti.Styles()
	s.Focused.Text = textStyle
	s.Blurred.Text = textStyle
	s.Cursor.Color = cursorColor
	s.Focused.Placeholder = focusedPlaceholderStyle
	s.Focused.Prompt = focusedPromptStyle
	s.Blurred.Placeholder = placeholderStyle
	ti.SetStyles(s)

	//textarea
	ta := textarea.New()
	ta.Placeholder = "Write your note here"
	ta.ShowLineNumbers = false
	ta.SetVirtualCursor(false)
	ta.SetStyles(textarea.DefaultDarkStyles())
	ta.Focus()

	//list
	noteList, err := listFile()
	if err != nil {
		noteList = []list.Item{}
	}
	finalList := list.New(noteList, list.NewDefaultDelegate(), 0, 0)
	finalList.Title = "All notes"
	finalList.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("254")).
		Padding(0, 1)

	return Model{
		NewFileInput: ti,
		State:        stateIdle,
		Style:        NewStyles(false),
		NoteTextArea: ta,
		List:         finalList,
	}
}
