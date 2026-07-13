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

	runningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	stoppedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	unknownStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220"))

	typeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("81"))
)

func (m Model) View() string {
	if m.loading {
		return m.loadingView()
	}

	if m.err != nil {
		return m.errorView()
	}

	switch m.screen {
	case resourcesScreen:
		return m.resourcesView()

	case environmentsScreen:
		return m.environmentsView()

	case resourceDetailsScreen:
		return m.resourceDetailsView()

	case deploymentsScreen:
		return m.deploymentsView()

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
				"%d/%d • j/k move • enter open • esc back • r refresh • q quit",
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

func (m Model) resourcesView() string {
	var view strings.Builder

	if m.project == nil ||
		len(m.project.Environments) == 0 {
		return "No environment selected."
	}

	environment := m.project.Environments[m.environmentCursor]

	title := fmt.Sprintf(
		"Coolify / %s / Environments / %s / Resources",
		m.project.Name,
		environment.Name,
	)

	view.WriteString(titleStyle.Render(title))
	view.WriteString("\n\n")

	if len(m.resources) == 0 {
		view.WriteString("No resources found.\n\n")
		view.WriteString(
			footerStyle.Render(
				"esc back • r refresh • q quit",
			),
		)
		return view.String()
	}

	start, end := m.visibleRange(
		m.resourceCursor,
		len(m.resources),
	)

	for index := start; index < end; index++ {
		resource := m.resources[index]

		cursor := "  "
		name := resource.Name

		if index == m.resourceCursor {
			cursor = "› "
			name = selectedStyle.Render(name)
		}

		view.WriteString(cursor)
		view.WriteString(name)
		view.WriteString(" ")
		view.WriteString(
			typeStyle.Render("[" + resource.Type + "]"),
		)
		view.WriteString(" ")
		view.WriteString(renderStatus(resource.Status))
		view.WriteString("\n")
	}

	view.WriteString("\n")
	view.WriteString(
		footerStyle.Render(
			fmt.Sprintf(
				"%d/%d • j/k move • enter details • esc back • r refresh • q quit",
				m.resourceCursor+1,
				len(m.resources),
			),
		),
	)

	return view.String()
}

func (m Model) resourceDetailsView() string {
	resource := m.selectedResource()
	if resource == nil {
		return "No resource selected."
	}

	var view strings.Builder

	title := fmt.Sprintf(
		"Coolify / %s / Environments / %s / Resources / %s",
		m.project.Name,
		m.project.Environments[m.environmentCursor].Name,
		resource.Name,
	)

	view.WriteString(titleStyle.Render(title))
	view.WriteString("\n\n")

	writeDetail := func(label string, value string) {
		view.WriteString(
			descriptionStyle.Render(label + ": "),
		)
		view.WriteString(value)
		view.WriteString("\n")
	}

	writeDetail("Name", resource.Name)
	writeDetail("Type", resource.Type)
	writeDetail("Status", renderStatus(resource.Status))
	writeDetail("UUID", resource.UUID)
	writeDetail(
		"Environment ID",
		fmt.Sprintf("%d", resource.EnvironmentID),
	)

	if resource.Description != nil &&
		*resource.Description != "" {
		writeDetail("Description", *resource.Description)
	}

	if resource.FQDN != nil &&
		*resource.FQDN != "" {
		writeDetail("FQDN", *resource.FQDN)
	}

	footer := "esc back • q quit"

	if strings.EqualFold(resource.Type, "application") {
		footer = "d deployments • esc back • q quit"
	}

	view.WriteString("\n")
	view.WriteString(footerStyle.Render(footer))

	return view.String()
}

func (m Model) deploymentsView() string {
	var view strings.Builder

	resource := m.selectedResource()
	if resource == nil {
		return "No application selected."
	}

	title := fmt.Sprintf(
		"Coolify / %s / Deployments",
		resource.Name,
	)

	view.WriteString(titleStyle.Render(title))
	view.WriteString("\n\n")

	if len(m.deployments) == 0 {
		view.WriteString("No deployments found.\n\n")
		view.WriteString(
			footerStyle.Render(
				"esc back • r refresh • q quit",
			),
		)
		return view.String()
	}

	start, end := m.visibleRange(
		m.deploymentCursor,
		len(m.deployments),
	)

	for index := start; index < end; index++ {
		deployment := m.deployments[index]

		cursor := "  "
		label := shortCommit(deployment.Commit)

		if deployment.CommitMessage != nil &&
			*deployment.CommitMessage != "" {
			label += " — " + singleLine(
				*deployment.CommitMessage,
			)
		}

		if index == m.deploymentCursor {
			cursor = "› "
			label = selectedStyle.Render(label)
		}

		view.WriteString(cursor)
		view.WriteString(renderDeploymentStatus(
			deployment.Status,
		))
		view.WriteString("  ")
		view.WriteString(label)
		view.WriteString("\n")
	}

	view.WriteString("\n")
	view.WriteString(
		footerStyle.Render(
			fmt.Sprintf(
				"%d/%d loaded • %d total • j/k move • esc back • r refresh • q quit",
				m.deploymentCursor+1,
				len(m.deployments),
				m.deploymentCount,
			),
		),
	)

	return view.String()
}

func renderDeploymentStatus(status string) string {
	lowerStatus := strings.ToLower(status)

	switch {
	case strings.Contains(lowerStatus, "finished"),
		strings.Contains(lowerStatus, "success"):
		return runningStyle.Render("● " + status)

	case strings.Contains(lowerStatus, "failed"),
		strings.Contains(lowerStatus, "error"):
		return stoppedStyle.Render("● " + status)

	case strings.Contains(lowerStatus, "progress"),
		strings.Contains(lowerStatus, "queued"),
		strings.Contains(lowerStatus, "running"):
		return unknownStyle.Render("● " + status)

	default:
		return descriptionStyle.Render("● " + status)
	}
}

func shortCommit(commit string) string {
	commit = strings.TrimSpace(commit)

	if len(commit) > 8 {
		return commit[:8]
	}

	if commit == "" {
		return "no commit"
	}

	return commit
}

func singleLine(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\r", " ")

	return strings.TrimSpace(value)
}

func renderStatus(status string) string {
	lowerStatus := strings.ToLower(status)

	switch {
	case strings.HasPrefix(lowerStatus, "running"):
		return runningStyle.Render("● " + status)

	case strings.Contains(lowerStatus, "exited"),
		strings.Contains(lowerStatus, "unhealthy"),
		strings.Contains(lowerStatus, "stopped"):
		return stoppedStyle.Render("● " + status)

	default:
		return unknownStyle.Render("● " + status)
	}
}

func (m Model) loadingView() string {
	title := "Coolify / Projects"

	if m.project != nil {
		switch m.screen {
		case environmentsScreen:
			title = fmt.Sprintf(
				"Coolify / %s / Environments",
				m.project.Name,
			)

		case resourcesScreen:
			if len(m.project.Environments) > 0 {
				environment := m.project.Environments[m.environmentCursor]
				title = fmt.Sprintf(
					"Coolify / %s / Environments / %s / Resources",
					m.project.Name,
					environment.Name,
				)
			}
		}
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
