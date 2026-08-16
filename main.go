package main

import (
	"fmt"
	"log"
	"os"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
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

	vaultDir = fmt.Sprintf("%s/.totion", homedir)
}

type model struct {
	newFileInput textinput.Model
	inputVisible bool
	currentFile  *os.File
	noteTextArea textarea.Model
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+n":
			m.inputVisible = true
		case "ctrl+s":
			if m.currentFile == nil {
				break
			}

			if err := m.currentFile.Truncate(0); err != nil {
				fmt.Println("can not save file :(")
				return m, nil
			}

			if _, err := m.currentFile.Seek(0, 0); err != nil {
				fmt.Println("can not save the file :(")
				return m, nil
			}

			if _, err := m.currentFile.WriteString(m.noteTextArea.Value()); err != nil {
				fmt.Println("can not save the file :(")
				return m, nil
			}

			if err := m.currentFile.Close(); err != nil {
				fmt.Println("can not close the file.")
			}

			m.currentFile = nil
			m.noteTextArea.SetValue("")

			return m, nil
		case "enter":
			if m.currentFile != nil {
				break	
			}
			//todo: creat folder
			fileName := m.newFileInput.Value()
			if fileName != "" {
				filepath := fmt.Sprintf("%s/%s.md", vaultDir, fileName)

				if _, err := os.Stat(filepath); err == nil {
					return m, nil
				}

				f, err := os.Create(filepath)
				if err != nil {
					log.Fatalf("%v", err)
				}

				m.currentFile = f
				m.inputVisible = false
				m.newFileInput.SetValue("")
			}
			return m, nil
		}
	}

	if m.inputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
	}

	return m, cmd
}

func (m model) View() tea.View {

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).
		PaddingLeft(2).
		PaddingRight(2)

	welcome := style.Render("Welcome to Totion 🧠")

	help := "Ctrl+N: new file . Ctrl+L: list . Esc: back/save . Ctrl+S: save . Ctrl+Q: quit"

	view := ""
	if m.inputVisible {
		view = m.newFileInput.View()
	}

	if m.currentFile != nil {
		view = m.noteTextArea.View()
	}

	prit := fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)

	return tea.NewView(prit)
}

func initializeModel() model {

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

	return model{
		newFileInput: ti,
		inputVisible: false,
		noteTextArea: ta,
	}
}

func main() {

	p := tea.NewProgram(initializeModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
