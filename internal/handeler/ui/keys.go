package ui

import "charm.land/bubbles/v2/key"

type keyMap struct {
	New     key.Binding
	List    key.Binding
	Rename  key.Binding
	Delete  key.Binding
	Preview key.Binding
	Save    key.Binding
	Back    key.Binding
	Help    key.Binding
	Quit    key.Binding
}

var keys = keyMap{
	New: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "new note"),
	),

	List: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "list note"),
	),

	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),

	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),

	Preview: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("ctrl+p", "preview"),
	),

	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),

	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),

	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.New, k.List, k.Back, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.New, k.List, k.Preview},
		{k.Rename, k.Delete, k.Save},
		{k.Back, k.Help, k.Quit},
	}
}
