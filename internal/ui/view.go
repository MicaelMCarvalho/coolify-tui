package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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
	if m.width < 80 || m.height < 24 {
		return fmt.Sprintf(
			"Terminal too small: %dx%d\nMinimum recommended size: 80x24",
			m.width,
			m.height,
		)
	}

	footerHeight := 1
	contentHeight := m.height - footerHeight

	leftWidth := max(30, m.width/3)
	rightWidth := m.width - leftWidth

	teamsHeight := 5
	environmentsHeight := 6

	remainingLeftHeight :=
		contentHeight -
			teamsHeight -
			environmentsHeight

	projectsHeight := remainingLeftHeight / 2
	resourcesHeight :=
		remainingLeftHeight - projectsHeight

	commandHeight := 5
	remainingRightHeight :=
		contentHeight - commandHeight

	detailsHeight := max(
		8,
		remainingRightHeight/4,
	)
	remainingListHeight :=
		remainingRightHeight - detailsHeight

	environmentVariableHeight :=
		remainingListHeight / 2

	deploymentsHeight :=
		remainingListHeight - environmentVariableHeight

	leftColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		m.teamsPane(
			leftWidth,
			teamsHeight,
		),
		m.projectsPane(
			leftWidth,
			projectsHeight,
		),
		m.environmentsPane(
			leftWidth,
			environmentsHeight,
		),
		m.resourcesPane(
			leftWidth,
			resourcesHeight,
		),
	)

	rightColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		m.detailsPane(
			rightWidth,
			detailsHeight,
		),
		m.environmentVariablesPane(
			rightWidth,
			environmentVariableHeight,
		),
		m.deploymentsPane(
			rightWidth,
			deploymentsHeight,
		),
		m.commandLogPane(
			rightWidth,
			commandHeight,
		),
	)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn,
		rightColumn,
	)

	footer := footerStyle.Render(
		"tab/shift+tab panel • 1-7 jump • j/k move • v reveal env • r refresh • q quit",
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		footer,
	)
}

func (m Model) renderPane(
	target panel,
	title string,
	body string,
	width int,
	height int,
) string {
	borderColor := lipgloss.Color("240")
	renderedTitle := descriptionStyle.Render(title)

	if m.activePanel == target {
		borderColor = lipgloss.Color("42")
		renderedTitle = selectedStyle.Render(title)
	}

	style := lipgloss.NewStyle().
		Width(max(1, width-2)).
		Height(max(1, height-2)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	return style.Render(
		renderedTitle + "\n" + body,
	)
}

func renderList(
	items []string,
	cursor int,
	maxRows int,
	maxWidth int,
) string {
	if len(items) == 0 {
		return descriptionStyle.Render("No items")
	}

	maxRows = max(1, maxRows)
	maxWidth = max(4, maxWidth)

	start := 0

	if cursor >= maxRows {
		start = cursor - maxRows + 1
	}

	end := min(
		len(items),
		start+maxRows,
	)

	var output strings.Builder

	for index := start; index < end; index++ {
		prefix := "  "
		item := ansi.Truncate(
			items[index],
			maxWidth-2,
			"…",
		)

		if index == cursor {
			prefix = "› "
			item = selectedStyle.Render(item)
		}

		output.WriteString(prefix)
		output.WriteString(item)
		output.WriteString("\n")
	}

	return strings.TrimRight(
		output.String(),
		"\n",
	)
}

func (m Model) teamsPane(
	width int,
	height int,
) string {
	items := make(
		[]string,
		0,
		len(m.teams),
	)

	for _, team := range m.teams {
		items = append(items, team.Name)
	}

	title := fmt.Sprintf(
		"[1] Teams (%d)",
		len(m.teams),
	)

	return m.renderPane(
		teamsPanel,
		title,
		renderList(
			items,
			m.teamCursor,
			height-3,
			width-4,
		),
		width,
		height,
	)
}

func (m Model) projectsPane(
	width int,
	height int,
) string {
	items := make(
		[]string,
		0,
		len(m.projects),
	)

	for _, project := range m.projects {
		items = append(items, project.Name)
	}

	title := fmt.Sprintf(
		"[2] Projects (%d)",
		len(m.projects),
	)

	return m.renderPane(
		projectsPanel,
		title,
		renderList(
			items,
			m.projectCursor,
			height-3,
			width-4,
		),
		width,
		height,
	)
}

func (m Model) environmentsPane(
	width int,
	height int,
) string {
	var items []string

	if m.project != nil {
		items = make(
			[]string,
			0,
			len(m.project.Environments),
		)

		for _, environment := range m.project.Environments {
			items = append(
				items,
				environment.Name,
			)
		}
	}

	title := fmt.Sprintf(
		"[3] Environments (%d)",
		len(items),
	)

	return m.renderPane(
		environmentsPanel,
		title,
		renderList(
			items,
			m.environmentCursor,
			height-3,
			width-4,
		),
		width,
		height,
	)
}

func (m Model) resourcesPane(
	width int,
	height int,
) string {
	items := make(
		[]string,
		0,
		len(m.resources),
	)

	for _, resource := range m.resources {
		resourceType := typeStyle.Render(
			"[" + resource.Type + "]",
		)

		item := fmt.Sprintf(
			"%s %s %s",
			resourceType,
			resource.Name,
			renderStatus(resource.Status),
		)

		items = append(items, item)
	}

	title := fmt.Sprintf(
		"[4] Resources (%d)",
		len(m.resources),
	)

	return m.renderPane(
		resourcesPanel,
		title,
		renderList(
			items,
			m.resourceCursor,
			height-3,
			width-4,
		),
		width,
		height,
	)
}

func (m Model) detailsPane(
	width int,
	height int,
) string {
	resource := m.selectedResource()

	if resource == nil {
		return m.renderPane(
			detailsPanel,
			"[5] Resource Details",
			descriptionStyle.Render(
				"Select a resource",
			),
			width,
			height,
		)
	}

	var body strings.Builder

	writeDetail := func(
		label string,
		value string,
	) {
		if strings.TrimSpace(value) == "" {
			return
		}

		body.WriteString(
			descriptionStyle.Render(
				label + ": ",
			),
		)
		body.WriteString(value)
		body.WriteString("\n")
	}

	writeDetail("Name", resource.Name)
	writeDetail("Type", resource.Type)
	writeDetail(
		"Status",
		renderStatus(resource.Status),
	)
	writeDetail("UUID", resource.UUID)
	writeDetail(
		"Environment ID",
		fmt.Sprintf(
			"%d",
			resource.EnvironmentID,
		),
	)

	if resource.Description != nil {
		writeDetail(
			"Description",
			*resource.Description,
		)
	}

	if resource.FQDN != nil {
		writeDetail("URL", *resource.FQDN)
	}

	return m.renderPane(
		detailsPanel,
		"[5] Resource Details",
		strings.TrimRight(
			body.String(),
			"\n",
		),
		width,
		height,
	)
}

func (m Model) environmentVariablesPane(
	width int,
	height int,
) string {
	resource := m.selectedResource()

	if resource == nil ||
		!strings.EqualFold(
			resource.Type,
			"application",
		) {
		return m.renderPane(
			environmentVariablesPanel,
			"[6] Environment Variables",
			descriptionStyle.Render(
				"Select an application",
			),
			width,
			height,
		)
	}

	items := make(
		[]string,
		0,
		len(m.environmentVariables),
	)

	for _, variable := range m.environmentVariables {
		value := "••••••••"

		if m.revealEnvironmentValues {
			value = variable.Value

			if value == "" {
				value = variable.RealValue
			}
			value = singleLine(value)

			if value == "" {
				if variable.IsShownOnce {
					value = "(hidden by Coolify)"
				} else {
					value = "(empty)"
				}

			}
		}

		flags := make([]string, 0, 3)

		if variable.IsBuildTime {
			flags = append(flags, "build")
		}

		if variable.IsCoolify {
			flags = append(flags, "coolify")
		}

		if variable.IsRuntime {
			flags = append(flags, "runtime")
		}

		if variable.IsPreview {
			flags = append(flags, "preview")
		}

		item := variable.Key + "=" + value

		if len(flags) > 0 {
			item += descriptionStyle.Render(
				" [" + strings.Join(flags, ", ") + "]",
			)
		}
		items = append(items, item)
	}

	title := fmt.Sprintf(
		"[6] Environment Variables (%d)",
		len(m.environmentVariables),
	)

	if m.revealEnvironmentValues {
		title += " (values revealed)"
	}

	return m.renderPane(
		environmentVariablesPanel,
		title,
		renderList(
			items,
			m.environmentVariablesCursor,
			height-3,
			width-4,
		),
		width,
		height,
	)
}

func (m Model) deploymentsPane(
	width int,
	height int,
) string {
	items := make(
		[]string,
		0,
		len(m.deployments),
	)

	for _, deployment := range m.deployments {
		item := fmt.Sprintf(
			"%s %s",
			renderDeploymentStatus(
				deployment.Status,
			),
			shortCommit(deployment.Commit),
		)

		if deployment.CommitMessage != nil &&
			*deployment.CommitMessage != "" {
			item += " " + singleLine(
				*deployment.CommitMessage,
			)
		}

		items = append(items, item)
	}

	title := fmt.Sprintf(
		"[7] Deployments (%d/%d)",
		len(m.deployments),
		m.deploymentCount,
	)

	return m.renderPane(
		deploymentsPanel,
		title,
		renderList(
			items,
			m.deploymentCursor,
			height-3,
			width-4,
		),
		width,
		height,
	)
}

func (m Model) commandLogPane(
	width int,
	height int,
) string {
	message := "Ready"

	if m.loading {
		message = unknownStyle.Render(
			"Loading…",
		)
	}

	if m.err != nil {
		message = errorStyle.Render(
			"Error: " + m.err.Error(),
		)
	}

	// panel(-1) means this informational panel
	// is never directly focused.
	return m.renderPane(
		panel(-1),
		"Command Log",
		message,
		width,
		height,
	)
}

func renderDeploymentStatus(
	status string,
) string {
	lowerStatus := strings.ToLower(status)

	switch {
	case strings.Contains(
		lowerStatus,
		"finished",
	),
		strings.Contains(
			lowerStatus,
			"success",
		):
		return runningStyle.Render(
			"● " + status,
		)

	case strings.Contains(
		lowerStatus,
		"failed",
	),
		strings.Contains(
			lowerStatus,
			"error",
		):
		return stoppedStyle.Render(
			"● " + status,
		)

	case strings.Contains(
		lowerStatus,
		"progress",
	),
		strings.Contains(
			lowerStatus,
			"queued",
		),
		strings.Contains(
			lowerStatus,
			"running",
		):
		return unknownStyle.Render(
			"● " + status,
		)

	default:
		return descriptionStyle.Render(
			"● " + status,
		)
	}
}

func renderStatus(status string) string {
	lowerStatus := strings.ToLower(status)

	switch {
	case strings.HasPrefix(
		lowerStatus,
		"running",
	):
		return runningStyle.Render(
			"● " + status,
		)

	case strings.Contains(
		lowerStatus,
		"exited",
	),
		strings.Contains(
			lowerStatus,
			"unhealthy",
		),
		strings.Contains(
			lowerStatus,
			"stopped",
		):
		return stoppedStyle.Render(
			"● " + status,
		)

	default:
		return unknownStyle.Render(
			"● " + status,
		)
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
	value = strings.ReplaceAll(
		value,
		"\n",
		" ",
	)
	value = strings.ReplaceAll(
		value,
		"\r",
		" ",
	)

	return strings.TrimSpace(value)
}
