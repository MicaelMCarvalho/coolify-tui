package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/micaelmcarvalho/coolify-tui/internal/coolify"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case teamsLoadedMsg:
		m.applyTeams(msg.teams)
		m.loading = false
		m.err = nil

	case projectsLoadedMsg:
		m.applyProjects(msg.projects)
		m.err = nil

		project := m.selectedProject()
		if project == nil {
			m.loading = false
			return m, nil
		}

		m.loading = true
		return m, m.loadProject(project.UUID)

	case projectLoadedMsg:
		selected := m.selectedProject()

		// Ignore stale responses caused by quick navigation.
		if selected == nil ||
			selected.UUID != msg.projectUUID {
			return m, nil
		}

		m.applyProject(msg.project)
		m.err = nil

		environment := m.selectedEnvironment()
		if environment == nil {
			m.loading = false
			return m, nil
		}

		m.loading = true
		return m, m.loadResources(environment.ID)

	case resourcesLoadedMsg:
		environment := m.selectedEnvironment()

		// Ignore a response for an environment that is no
		// longer selected.
		if environment == nil ||
			environment.ID != msg.environmentID {
			return m, nil
		}

		m.applyResources(msg.resources)
		m.err = nil

		resource := m.selectedResource()
		if resource == nil ||
			!strings.EqualFold(
				resource.Type,
				"application",
			) {
			m.loading = false
			return m, nil
		}

		m.loading = true
		m.deploymentSkip = 0

		return m, tea.Batch(
			m.loadEnvironmentVariables(resource.UUID),
			m.loadDeployments(resource.UUID, 0),
		)

	case deploymentsLoadedMsg:
		resource := m.selectedResource()

		// Ignore an older deployment response if the user
		// has selected a different resource.
		if resource == nil ||
			resource.UUID != msg.resourceUUID {
			return m, nil
		}

		m.deployments = msg.result.Deployments
		m.deploymentCount = msg.result.Count
		m.deploymentSkip = msg.skip
		m.deploymentCursor = 0
		m.loading = false
		m.err = nil

	case deploymentDetailsLoadedMsg:
		selected := m.selectedDeployment()

		if !m.deploymentDetailsOpen ||
			selected == nil ||
			selected.DeploymentUUID != msg.deployment.DeploymentUUID {
			return m, nil
		}

		details := msg.deployment
		m.deploymentDetails = &details
		m.deploymentLogOffset = 0
		m.loading = false
		m.err = nil

	case environmentVariablesLoadedMsg:
		resource := m.selectedResource()

		if resource == nil ||
			resource.UUID != msg.resourceUUID {
			return m, nil
		}

		m.environmentVariables = msg.variables
		m.environmentVariablesCursor = 0
		m.loading = false
		m.err = nil

	case errMsg:
		m.loading = false
		m.err = msg.err

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(
	msg tea.KeyMsg,
) (tea.Model, tea.Cmd) {
	if m.deploymentDetailsOpen {
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.deploymentDetailsOpen = false
			m.deploymentDetails = nil
			m.deploymentLogOffset = 0
			m.loading = false
			m.err = nil

			return m, nil

		case "up", "k":
			if m.deploymentLogOffset > 0 {
				m.deploymentLogOffset--
			}

		case "down", "j":
			if m.deploymentDetails != nil {
				lines := deploymentLogLines(
					m.deploymentDetails.Logs,
				)

				if m.deploymentLogOffset <
					max(0, len(lines)-1) {
					m.deploymentLogOffset++
				}
			}

		case "g", "home":
			m.deploymentLogOffset = 0

		case "G", "end":
			if m.deploymentDetails != nil {
				lines := deploymentLogLines(
					m.deploymentDetails.Logs,
				)

				m.deploymentLogOffset =
					max(0, len(lines)-1)
			}

		case "r":
			selected := m.selectedDeployment()
			if selected == nil {
				return m, nil
			}

			m.loading = true
			m.err = nil

			return m, m.loadDeploymentDetails(
				selected.DeploymentUUID,
			)
		}

		return m, nil
	}

	if m.filtering {
		return m.handleFilterKey(msg)
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "tab":
		m.activePanel =
			(m.activePanel + 1) % 7

	case "shift+tab", "backtab":
		m.activePanel =
			(m.activePanel + 6) % 7

	case "1":
		m.activePanel = teamsPanel

	case "2":
		m.activePanel = projectsPanel

	case "3":
		m.activePanel = environmentsPanel

	case "4":
		m.activePanel = resourcesPanel

	case "5":
		m.activePanel = detailsPanel

	case "6":
		m.activePanel = environmentVariablesPanel

	case "7":
		m.activePanel = deploymentsPanel

	case "enter":
		if m.activePanel == deploymentsPanel {
			deployment := m.selectedDeployment()
			if deployment == nil {
				return m, nil
			}

			m.deploymentDetailsOpen = true
			m.deploymentDetails = nil
			m.deploymentLogOffset = 0
			m.loading = true
			m.err = nil

			return m, m.loadDeploymentDetails(
				deployment.DeploymentUUID,
			)
		}

		m.focusNextPanel()

	case "esc":
		if filterablePanel(m.activePanel) &&
			m.filters[m.activePanel] != "" {
			m.filters[m.activePanel] = ""
			return m, m.moveToBoundary(true)
		}

		m.focusPreviousPanel()

	case "up", "k":
		cmd := m.moveCursor(-1)
		return m, cmd

	case "down", "j":
		cmd := m.moveCursor(1)
		return m, cmd

	case "g", "home":
		cmd := m.moveToBoundary(true)
		return m, cmd

	case "G", "end":
		cmd := m.moveToBoundary(false)
		return m, cmd

	case "r":
		cmd := m.refreshActivePanel()
		return m, cmd

	case "n":
		cmd := m.nextDeploymentPage()
		return m, cmd

	case "p":
		cmd := m.previousDeploymentPage()
		return m, cmd

	case "v":
		if len(m.environmentVariables) > 0 {
			m.revealEnvironmentValues =
				!m.revealEnvironmentValues
		}

	case "/":
		if !filterablePanel(m.activePanel) {
			return m, nil
		}

		m.filtering = true
		m.filterPanel = m.activePanel
		m.filterOriginal = m.filters[m.activePanel]
		m.filterInput = m.filters[m.activePanel]
		return m, nil
	}

	return m, nil
}

func (m *Model) focusNextPanel() {
	if m.activePanel < deploymentsPanel {
		m.activePanel++
	}
}

func (m *Model) focusPreviousPanel() {
	if m.activePanel > teamsPanel {
		m.activePanel--
	}
}

func (m *Model) moveCursor(
	change int,
) tea.Cmd {
	switch m.activePanel {
	case teamsPanel:
		next := m.teamCursor + change

		if next >= 0 && next < len(m.teams) {
			m.teamCursor = next
		}

		// A Coolify API token is scoped to one team, so
		// changing this cursor cannot switch API context yet.
		return nil

	case projectsPanel:
		indices := m.filteredIndices(
			projectsPanel,
		)

		next, ok := nextFilteredIndex(
			indices,
			m.projectCursor,
			change,
		)

		if !ok || next == m.projectCursor {
			return nil
		}

		m.projectCursor = next
		m.clearAfterProject()
		m.loading = true
		m.err = nil

		project := m.selectedProject()
		if project == nil {
			m.loading = false
			return nil
		}

		return m.loadProject(project.UUID)

	case environmentsPanel:
		if m.project == nil {
			return nil
		}

		indices := m.filteredIndices(
			environmentsPanel,
		)

		next, ok := nextFilteredIndex(
			indices,
			m.environmentCursor,
			change,
		)

		if !ok || next == m.environmentCursor {
			return nil
		}

		m.environmentCursor = next
		m.clearAfterEnvironment()
		m.loading = true
		m.err = nil

		environment := m.selectedEnvironment()
		if environment == nil {
			m.loading = false
			return nil
		}

		return m.loadResources(environment.ID)

	case resourcesPanel:
		indices := m.filteredIndices(
			resourcesPanel,
		)

		next, ok := nextFilteredIndex(
			indices,
			m.resourceCursor,
			change,
		)

		if !ok || next == m.resourceCursor {
			return nil
		}

		m.resourceCursor = next
		m.clearAfterResource()
		m.err = nil

		resource := m.selectedResource()
		if resource == nil ||
			!strings.EqualFold(
				resource.Type,
				"application",
			) {
			m.loading = false
			return nil
		}

		m.loading = true

		return tea.Batch(
			m.loadEnvironmentVariables(resource.UUID),
			m.loadDeployments(resource.UUID, 0),
		)

	case environmentVariablesPanel:
		indices := m.filteredIndices(
			environmentVariablesPanel,
		)

		next, ok := nextFilteredIndex(
			indices,
			m.environmentVariablesCursor,
			change,
		)

		if !ok || next == m.environmentVariablesCursor {
			return nil
		}

		m.environmentVariablesCursor = next

	case deploymentsPanel:
		next := m.deploymentCursor + change

		if next >= 0 &&
			next < len(m.deployments) {
			m.deploymentCursor = next
		}
	}

	return nil
}

func (m *Model) moveToBoundary(
	first bool,
) tea.Cmd {
	switch m.activePanel {
	case teamsPanel:
		if len(m.teams) == 0 {
			return nil
		}

		if first {
			m.teamCursor = 0
		} else {
			m.teamCursor = len(m.teams) - 1
		}

	case projectsPanel:
		indices := m.filteredIndices(
			projectsPanel,
		)

		if len(indices) == 0 {
			return nil
		}

		target := indices[0]
		if !first {
			target = indices[len(indices)-1]
		}

		if target == m.projectCursor {
			return nil
		}

		m.projectCursor = target
		m.clearAfterProject()
		m.loading = true
		m.err = nil

		project := m.selectedProject()
		if project == nil {
			m.loading = false
			return nil
		}

		return m.loadProject(project.UUID)

	case environmentsPanel:
		if m.project == nil {
			return nil
		}

		indices := m.filteredIndices(
			environmentsPanel,
		)

		if len(indices) == 0 {
			return nil
		}

		target := indices[0]

		if !first {
			target =
				indices[len(indices)-1]
		}

		if target == m.environmentCursor {
			return nil
		}

		m.environmentCursor = target
		m.clearAfterEnvironment()
		m.loading = true
		m.err = nil

		environment := m.selectedEnvironment()
		if environment == nil {
			m.loading = false
			return nil
		}

		return m.loadResources(environment.ID)

	case resourcesPanel:
		indices := m.filteredIndices(
			resourcesPanel,
		)
		if len(indices) == 0 {
			return nil
		}

		target := indices[0]

		if !first {
			target = indices[len(indices)-1]
		}

		if target == m.resourceCursor {
			return nil
		}

		m.resourceCursor = target
		m.clearAfterResource()
		m.err = nil

		resource := m.selectedResource()
		if resource == nil ||
			!strings.EqualFold(
				resource.Type,
				"application",
			) {
			m.loading = false
			return nil
		}

		m.loading = true

		return tea.Batch(
			m.loadEnvironmentVariables(resource.UUID),
			m.loadDeployments(resource.UUID, 0),
		)

	case environmentVariablesPanel:
		indices := m.filteredIndices(
			environmentVariablesPanel,
		)

		if len(indices) == 0 {
			return nil
		}

		if first {
			m.environmentVariablesCursor =
				indices[0]
		} else {
			m.environmentVariablesCursor =
				indices[len(indices)-1]
		}

	case deploymentsPanel:
		if len(m.deployments) == 0 {
			return nil
		}

		if first {
			m.deploymentCursor = 0
		} else {
			m.deploymentCursor =
				len(m.deployments) - 1
		}
	}

	return nil
}

func (m *Model) refreshActivePanel() tea.Cmd {
	m.err = nil
	m.loading = true

	switch m.activePanel {
	case teamsPanel:
		return m.loadTeams()

	case projectsPanel:
		return m.loadProjects()

	case environmentsPanel:
		project := m.selectedProject()
		if project == nil {
			m.loading = false
			return nil
		}

		return m.loadProject(project.UUID)

	case resourcesPanel, detailsPanel:
		environment := m.selectedEnvironment()
		if environment == nil {
			m.loading = false
			return nil
		}

		return m.loadResources(environment.ID)

	case environmentVariablesPanel:
		resource := m.selectedResource()

		if resource == nil ||
			!strings.EqualFold(
				resource.Type,
				"application",
			) {
			m.loading = false
			return nil
		}

		return m.loadEnvironmentVariables(
			resource.UUID,
		)

	case deploymentsPanel:
		resource := m.selectedResource()
		if resource == nil ||
			!strings.EqualFold(
				resource.Type,
				"application",
			) {
			m.loading = false
			return nil
		}

		return m.loadDeployments(
			resource.UUID,
			m.deploymentSkip,
		)
	}

	m.loading = false
	return nil
}

func (m *Model) nextDeploymentPage() tea.Cmd {
	if m.activePanel != deploymentsPanel ||
		m.loading {
		return nil
	}

	nextSkip :=
		m.deploymentSkip + m.deploymentTake

	if nextSkip >= m.deploymentCount {
		return nil
	}

	resource := m.selectedResource()
	if resource == nil ||
		!strings.EqualFold(
			resource.Type,
			"application",
		) {
		return nil
	}

	m.loading = true
	m.err = nil

	return m.loadDeployments(
		resource.UUID,
		nextSkip,
	)
}

func (m *Model) previousDeploymentPage() tea.Cmd {
	if m.activePanel != deploymentsPanel ||
		m.loading ||
		m.deploymentSkip == 0 {
		return nil
	}

	resource := m.selectedResource()
	if resource == nil ||
		!strings.EqualFold(
			resource.Type,
			"application",
		) {
		return nil
	}

	previousSkip := max(
		0,
		m.deploymentSkip-m.deploymentTake,
	)

	m.loading = true
	m.err = nil

	return m.loadDeployments(
		resource.UUID,
		previousSkip,
	)
}

func (m *Model) applyTeams(
	teams []coolify.Team,
) {
	selectedID := -1

	if selected := m.selectedTeam(); selected != nil {
		selectedID = selected.ID
	}

	m.teams = teams
	m.teamCursor = 0

	for index, team := range m.teams {
		if team.ID == selectedID {
			m.teamCursor = index
			break
		}
	}
}

func (m *Model) applyProjects(
	projects []coolify.Project,
) {
	selectedUUID := ""

	if selected := m.selectedProject(); selected != nil {
		selectedUUID = selected.UUID
	}

	m.projects = projects
	m.projectCursor = 0

	for index, project := range m.projects {
		if project.UUID == selectedUUID {
			m.projectCursor = index
			break
		}
	}

	m.clearAfterProject()
}

func (m *Model) applyProject(
	project coolify.ProjectDetails,
) {
	selectedEnvironmentID := -1

	if selected := m.selectedEnvironment(); selected != nil {
		selectedEnvironmentID = selected.ID
	}

	projectCopy := project
	m.project = &projectCopy
	m.environmentCursor = 0

	for index, environment := range m.project.Environments {
		if environment.ID == selectedEnvironmentID {
			m.environmentCursor = index
			break
		}
	}

	m.clearAfterEnvironment()
}

func (m *Model) applyResources(
	resources []coolify.Resource,
) {
	selectedUUID := ""

	if selected := m.selectedResource(); selected != nil {
		selectedUUID = selected.UUID
	}

	m.resources = resources
	m.resourceCursor = 0

	for index, resource := range m.resources {
		if resource.UUID == selectedUUID {
			m.resourceCursor = index
			break
		}
	}

	m.clearAfterResource()
}

func (m *Model) clearAfterProject() {
	m.project = nil
	m.environmentCursor = 0
	m.resources = nil
	m.resourceCursor = 0
	m.clearAfterResource()
}

func (m *Model) clearAfterEnvironment() {
	m.resources = nil
	m.resourceCursor = 0
	m.clearAfterResource()
}

func (m *Model) clearAfterResource() {
	m.environmentVariables = nil
	m.environmentVariablesCursor = 0
	m.revealEnvironmentValues = false

	m.deployments = nil
	m.deploymentCount = 0
	m.deploymentCursor = 0
	m.deploymentSkip = 0
}
