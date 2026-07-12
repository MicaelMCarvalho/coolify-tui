package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	descriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

func (m Model) View() string {
	var view strings.Builder

	view.WriteString(titleStyle.Render("Coolify Projects"))
	view.WriteString("\n\n")

	if m.loading {
		view.WriteString("Loading projects…\n")
		return view.String()
	}

	if m.err != nil {
		view.WriteString(errorStyle.Render(
			fmt.Sprintf("Error: %v", m.err),
		))
		view.WriteString("\n\n")
		view.WriteString(footerStyle.Render("r retry • q quit"))
		return view.String()
	}

	if len(m.projects) == 0 {
		view.WriteString("No projects found.\n\n")
		view.WriteString(footerStyle.Render("r refresh • q quit"))
		return view.String()
	}

	start, end := m.visibleRange()

	for index := start; index < end; index++ {
		project := m.projects[index]

		cursor := "  "
		name := project.Name

		if index == m.cursor {
			cursor = "› "
			name = selectedStyle.Render(name)
		}

		view.WriteString(cursor)
		view.WriteString(name)

		if project.Description != "" {
			view.WriteString(
				descriptionStyle.Render(
					" — " + project.Description,
				),
			)
		}

		view.WriteString("\n")
	}

	view.WriteString("\n")
	view.WriteString(
		footerStyle.Render(
			fmt.Sprintf(
				"%d/%d • j/k move • g/G first/last • r refresh • q quit",
				m.cursor+1,
				len(m.projects),
			),
		),
	)

	return view.String()
}

func (m Model) visibleRange() (int, int) {
	available := m.height - 6
	if available < 1 {
		available = 10
	}

	start := 0
	if m.cursor >= available {
		start = m.cursor - available + 1
	}

	end := min(len(m.projects), start+available)

	return start, end
}
