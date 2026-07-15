package ui

import (
	"encoding/json"
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

	if m.helpOpen {
		backgroundModel := m
		backgroundModel.helpOpen = false

		background := backgroundModel.View()
		popup := m.helpPopup()

		return overlayCentered(
			background,
			popup,
			m.width,
			m.height,
		)
	}

	if m.deploymentDetailsOpen {
		return m.deploymentDetailsView()
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

	footerText :=
		"tab/shift+tab panel • 1-7 jump • j/k move • / filter • ? help • r refresh • q quit"

	if m.filtering {
		footerText = fmt.Sprintf(
			"Filter: %s█ • enter accept • esc cancel • ctrl+u clear",
			m.filterInput,
		)
	} else if filterablePanel(m.activePanel) &&
		m.filters[m.activePanel] != "" {
		footerText = fmt.Sprintf(
			"Filter: %s • / edit • esc clear • j/k move",
			m.filters[m.activePanel],
		)
	}

	footer := footerStyle.Render(footerText)

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

	indices := m.filteredIndices(projectsPanel)

	items = filterItemsByIndices(
		items,
		indices,
	)

	cursor := filteredCursorPosition(
		indices,
		m.projectCursor,
	)

	title := fmt.Sprintf(
		"[2] Projects (%d/%d)",
		len(items),
		len(m.projects),
	)

	return m.renderPane(
		projectsPanel,
		title,
		renderList(
			items,
			cursor,
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
	total := 0

	if m.project != nil {
		total = len(m.project.Environments)

		items = make(
			[]string,
			0,
			total,
		)

		for _, environment := range m.project.Environments {
			items = append(
				items,
				environment.Name,
			)
		}
	}

	indices := m.filteredIndices(
		environmentsPanel,
	)

	items = filterItemsByIndices(
		items,
		indices,
	)

	cursor := filteredCursorPosition(
		indices,
		m.environmentCursor,
	)

	title := fmt.Sprintf(
		"[3] Environments (%d/%d)",
		len(items),
		total,
	)

	return m.renderPane(
		environmentsPanel,
		title,
		renderList(
			items,
			cursor,
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

	indices := m.filteredIndices(
		resourcesPanel,
	)

	items = filterItemsByIndices(
		items,
		indices,
	)

	cursor := filteredCursorPosition(
		indices,
		m.resourceCursor,
	)

	title := fmt.Sprintf(
		"[4] Resources (%d/%d)",
		len(items),
		len(m.resources),
	)

	return m.renderPane(
		resourcesPanel,
		title,
		renderList(
			items,
			cursor,
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

	indices := m.filteredIndices(
		environmentVariablesPanel,
	)

	items = filterItemsByIndices(
		items,
		indices,
	)

	cursor := filteredCursorPosition(
		indices,
		m.environmentVariablesCursor,
	)

	title := fmt.Sprintf(
		"[6] Environment Variables (%d/%d)",
		len(items),
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
			cursor,
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

	indices := m.filteredIndices(
		deploymentsPanel,
	)

	items = filterItemsByIndices(
		items,
		indices,
	)

	cursor := filteredCursorPosition(
		indices,
		m.deploymentCursor,
	)

	title := fmt.Sprintf(
		"[7] Deployments (%d/%d shown, %d total)",
		len(items),
		len(m.deployments),
		m.deploymentCount,
	)

	return m.renderPane(
		deploymentsPanel,
		title,
		renderList(
			items,
			cursor,
			height-3,
			width-4,
		),
		width,
		height,
	)
}

func (m Model) deploymentDetailsView() string {
	footer := footerStyle.Render(
		"j/k scroll logs • g/G top/bottom • r refresh • esc back • q quit",
	)

	paneHeight := m.height - 1

	if m.err != nil {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderPane(
				deploymentsPanel,
				"[7] Deployment Details",
				errorStyle.Render(
					"Error: "+m.err.Error(),
				),
				m.width,
				paneHeight,
			),
			footer,
		)
	}

	if m.deploymentDetails == nil {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderPane(
				deploymentsPanel,
				"[7] Deployment Details",
				unknownStyle.Render("Loading…"),
				m.width,
				paneHeight,
			),
			footer,
		)
	}

	deployment := m.deploymentDetails

	commitMessage := ""
	if deployment.CommitMessage != nil {
		commitMessage = singleLine(
			*deployment.CommitMessage,
		)
	}

	finishedAt := ""
	if deployment.FinishedAt != nil {
		finishedAt = *deployment.FinishedAt
	}

	bodyLines := []string{
		"Status: " + renderDeploymentStatus(
			deployment.Status,
		),
		"Application: " + deployment.ApplicationName,
		"Deployment UUID: " +
			deployment.DeploymentUUID,
		"Commit: " + deployment.Commit,
		"Commit message: " + commitMessage,
		"Server: " + deployment.ServerName,
		"Created: " + deployment.CreatedAt,
		"Updated: " + deployment.UpdatedAt,
		"Finished: " + finishedAt,
		"",
		titleStyle.Render("Logs"),
	}

	// renderPane reserves one row for its title and two
	// rows for its border.
	availableLogRows :=
		deploymentLogPageSize(
			m.height,
		)

	logLines := deploymentLogLines(
		deployment.Logs,
	)

	maxStart := max(
		0,
		len(logLines)-availableLogRows,
	)

	start := min(
		m.deploymentLogOffset,
		maxStart,
	)

	end := min(
		len(logLines),
		start+availableLogRows,
	)

	logHeading := "Logs (0 of 0)"

	if len(logLines) > 0 {
		logHeading = fmt.Sprintf(
			"Logs (%d-%d of %d)",
			start+1,
			end,
			len(logLines),
		)
	}

	bodyLines[len(bodyLines)-1] =
		titleStyle.Render(logHeading)

	for _, line := range logLines[start:end] {
		bodyLines = append(
			bodyLines,
			ansi.Truncate(
				line,
				max(4, m.width-6),
				"…",
			),
		)
	}

	body := strings.Join(bodyLines, "\n")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderPane(
			deploymentsPanel,
			"[7] Deployment Details",
			body,
			m.width,
			paneHeight,
		),
		footer,
	)
}

func deploymentLogLines(
	raw json.RawMessage,
) []string {
	if len(raw) == 0 ||
		string(raw) == "null" {
		return []string{"No logs available"}
	}

	// Coolify can return logs as an encoded JSON string.
	var encoded string
	if err := json.Unmarshal(raw, &encoded); err == nil {
		encoded = strings.TrimSpace(encoded)

		if encoded == "" {
			return []string{"No logs available"}
		}

		if json.Valid([]byte(encoded)) {
			raw = json.RawMessage(encoded)
		} else {
			return splitLogLines(encoded)
		}
	}

	var entries []struct {
		Command   string `json:"command"`
		Output    string `json:"output"`
		Timestamp string `json:"timestamp"`
		Hidden    bool   `json:"hidden"`
	}

	if err := json.Unmarshal(raw, &entries); err == nil {
		lines := make([]string, 0)

		for _, entry := range entries {
			// Avoid showing Coolify commands explicitly marked
			// as hidden because they may contain secrets.
			if entry.Hidden {
				continue
			}

			prefix := ""
			if entry.Timestamp != "" {
				prefix = entry.Timestamp + " "
			}

			outputLines := splitLogLines(entry.Output)

			for _, line := range outputLines {
				lines = append(
					lines,
					prefix+line,
				)
			}
		}

		if len(lines) == 0 {
			return []string{"No visible logs available"}
		}

		return lines
	}

	return splitLogLines(string(raw))
}

func (m Model) helpPopup() string {
	leftLines := []string{
		titleStyle.Render("Navigation"),
		"tab          next panel",
		"shift+tab    previous panel",
		"1-7          focus panel",
		"j / k        move",
		"g / G        first / last",
		"enter        open / next",
		"esc          back",
		"",
		titleStyle.Render("Filtering"),
		"/            filter panel",
		"enter        accept",
		"esc          cancel / clear",
		"ctrl+u       clear input",
	}

	rightLines := []string{
		titleStyle.Render("Resources"),
		"v            reveal env values",
		"r            refresh",
		"",
		titleStyle.Render("Deployments"),
		"n / p        next / previous page",
		"enter        details and logs",
		"j / k        scroll logs",
		"g / G        top / bottom",
		"",
		titleStyle.Render("Global"),
		"?            close help",
		"q / ctrl+c   quit",
	}

	popupWidth := min(
		76,
		max(50, m.width-8),
	)

	columnWidth := max(
		22,
		(popupWidth-6)/2,
	)

	left := lipgloss.NewStyle().
		Width(columnWidth).
		Render(
			strings.Join(leftLines, "\n"),
		)

	right := lipgloss.NewStyle().
		Width(columnWidth).
		Render(
			strings.Join(rightLines, "\n"),
		)

	columns := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		right,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		selectedStyle.Render("[?] Keyboard Help"),
		"",
		columns,
		"",
		footerStyle.Render(
			"? or esc close",
		),
	)

	return lipgloss.NewStyle().
		Width(popupWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("42")).
		Background(lipgloss.Color("235")).
		Padding(0, 1).
		Render(content)
}

func deploymentLogPageSize(
	height int,
) int {
	paneHeight := height - 1
	maxBodyRows := max(1, paneHeight-3)

	const deploymentDetailsRows = 11

	return max(
		1,
		maxBodyRows-deploymentDetailsRows,
	)
}

func splitLogLines(value string) []string {
	value = strings.ReplaceAll(
		value,
		"\r\n",
		"\n",
	)

	value = strings.TrimSpace(value)

	if value == "" {
		return []string{"No logs available"}
	}

	return strings.Split(value, "\n")
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

func overlayCentered(
	background string,
	foreground string,
	width int,
	height int,
) string {
	backgroundLines := strings.Split(
		background,
		"\n",
	)

	foregroundLines := strings.Split(
		foreground,
		"\n",
	)

	foregroundWidth := lipgloss.Width(
		foreground,
	)

	foregroundHeight := lipgloss.Height(
		foreground,
	)

	x := max(
		0,
		(width-foregroundWidth)/2,
	)

	y := max(
		0,
		(height-foregroundHeight)/2,
	)

	// Make sure the background has exactly enough rows.
	for len(backgroundLines) < height {
		backgroundLines = append(
			backgroundLines,
			"",
		)
	}

	for row, foregroundLine := range foregroundLines {
		targetRow := y + row

		if targetRow < 0 ||
			targetRow >= len(backgroundLines) {
			continue
		}

		backgroundLine := padANSILine(
			backgroundLines[targetRow],
			width,
		)

		popupLine := padANSILine(
			foregroundLine,
			foregroundWidth,
		)

		left := ansi.Cut(
			backgroundLine,
			0,
			x,
		)

		right := ansi.Cut(
			backgroundLine,
			x+foregroundWidth,
			width,
		)

		backgroundLines[targetRow] =
			left + popupLine + right
	}

	if len(backgroundLines) > height {
		backgroundLines =
			backgroundLines[:height]
	}

	return strings.Join(
		backgroundLines,
		"\n",
	)
}

func padANSILine(
	line string,
	width int,
) string {
	line = ansi.Truncate(
		line,
		width,
		"",
	)

	missing := width -
		ansi.StringWidth(line)

	if missing > 0 {
		line += strings.Repeat(
			" ",
			missing,
		)
	}

	return line
}
