package ui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/Ashoke15/notes-app/internal/vault"
)

const autoSaveDelay = 1500 * time.Millisecond

type autoSaveMsg struct {
	gen int
}

func scheduleAutosave(gen int) tea.Cmd {
	return tea.Tick(autoSaveDelay, func(t time.Time) tea.Msg {
		return autoSaveMsg{gen: gen}
	})
}

func (m Model) handleAutosave(msg autoSaveMsg) (tea.Model, tea.Cmd) {
	if msg.gen != m.autosavegen {
		return m, nil
	}

	if m.CurrentFile == nil || !m.Dirty {
		return m, nil
	}

	if err := vault.Write(m.CurrentFile, m.NoteTextArea.Value()); err != nil {
		m.StatusMsg = "Autosave failed"
		return m, nil
	}

	m.Dirty = false
	m.StatusMsg = "Autosave"
	return m, nil
}
