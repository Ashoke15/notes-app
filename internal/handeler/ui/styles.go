package ui

import (
	"charm.land/lipgloss/v2"
)

func NewStyles(darkBG bool) Styles {
	lightDark := lipgloss.LightDark(darkBG)

	return Styles{
		app: lipgloss.NewStyle().
			Padding(1, 2),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1),
		statusMessage: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#04B575"), lipgloss.Color("#04B575"))),
		danger: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")),
		confirmBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF5555")).
			Padding(1, 3),
		confirmHint: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		keyHint: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5555")).
			Padding(0, 1),
		renameBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			Padding(1, 2),
		noteBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			Padding(0, 1),
	}
}
