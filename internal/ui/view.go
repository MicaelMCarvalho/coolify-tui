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
	if m.loading {
		return m.loadingView()
	}

	if m.err != nil {
		return m.errorView()
	}

	switch m.screen {
	case environmentsScreen:
		return m.environmentsView()

	default:
		return m.projectsView()
	}
}

func (m Model) projectsView() string {
	var view strings.Builder

	view.WriteString(titleStyle.Render("Coolify / Projects"))
	view.WriteString("\n\n")

	if len(m.projects) == 0 {
		view.WriteString("No projects found.\n\n")
		view.WriteString(
			footerStyle.Render("r refresh • q quit"),
		)
		return view.String()
	}

	start, end := m.visibleRange(
		m.projectCursor,
		len(m.projects),
	)

	for index := start; index < end; index++ {
		project := m.projects[index]

		cursor := "  "
		name := project.Name

		if index == m.projectCursor {
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
				"%d/%d • j/k move • enter open • r refresh • q quit",
				m.projectCursor+1,
				len(m.projects),
			),
		),
	)

	return view.String()
}

func (m Model) environmentsView() string {
	var view strings.Builder

	if m.project == nil {
		return "No project selected."
	}

	title := fmt.Sprintf(
		"Coolify / %s / Environments",
		m.project.Name,
	)

	view.WriteString(titleStyle.Render(title))
	view.WriteString("\n\n")

	environments := m.project.Environments

	if len(environments) == 0 {
		view.WriteString("No environments found.\n\n")
		view.WriteString(
			footerStyle.Render(
				"esc back • r refresh • q quit",
			),
		)
		return view.String()
	}

	start, end := m.visibleRange(
		m.environmentCursor,
		len(environments),
	)

	for index := start; index < end; index++ {
		environment := environments[index]

		cursor := "  "
		name := environment.Name

		if index == m.environmentCursor {
			cursor = "› "
			name = selectedStyle.Render(name)
		}

		view.WriteString(cursor)
		view.WriteString(name)

		if environment.Description != nil &&
			*environment.Description != "" {
			view.WriteString(
				descriptionStyle.Render(
					" — " + *environment.Description,
				),
			)
		}

		view.WriteString("\n")
	}

	view.WriteString("\n")
	view.WriteString(
		footerStyle.Render(
			fmt.Sprintf(
				"%d/%d • j/k move • esc back • r refresh • q quit",
				m.environmentCursor+1,
				len(environments),
			),
		),
	)

	return view.String()
}

func (m Model) loadingView() string {
	title := "Coolify / Projects"

	if m.screen == environmentsScreen &&
		m.project != nil {
		title = fmt.Sprintf(
			"Coolify / %s / Environments",
			m.project.Name,
		)
	}

	return titleStyle.Render(title) +
		"\n\nLoading…\n\n" +
		footerStyle.Render("q quit")
}

func (m Model) errorView() string {
	return titleStyle.Render("Coolify") +
		"\n\n" +
		errorStyle.Render(fmt.Sprintf("Error: %v", m.err)) +
		"\n\n" +
		footerStyle.Render("r retry • esc back • q quit")
}

func (m Model) visibleRange(
	cursor int,
	total int,
) (int, int) {
	available := m.height - 6
	if available < 1 {
		available = 10
	}

	start := 0
	if cursor >= available {
		start = cursor - available + 1
	}

	end := min(total, start+available)

	return start, end
}
