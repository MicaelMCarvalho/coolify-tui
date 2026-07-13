package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case projectsLoadedMsg:
		m.projects = msg.projects
		m.loading = false
		m.err = nil

		if m.projectCursor >= len(m.projects) {
			m.projectCursor = max(0, len(m.projects)-1)
		}

	case projectLoadedMsg:
		project := msg.project
		m.project = &project
		m.screen = environmentsScreen
		m.environmentCursor = 0
		m.loading = false
		m.err = nil

	case resourcesLoadedMsg:
		m.resources = msg.resources
		m.resourceCursor = 0
		m.screen = resourcesScreen
		m.loading = false
		m.err = nil

	case deploymentsLoadedMsg:
		m.deployments = msg.result.Deployments
		m.deploymentCount = msg.result.Count
		m.deploymentSkip = msg.skip
		m.deploymentCursor = 0
		m.screen = deploymentsScreen
		m.loading = false
		m.err = nil

	case errMsg:
		m.loading = false
		m.err = msg.err

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "esc":
			switch m.screen {
			case resourcesScreen:
				m.screen = environmentsScreen
				m.resources = nil
				m.resourceCursor = 0
				m.loading = false
				m.err = nil

			case environmentsScreen:
				m.screen = projectsScreen
				m.project = nil
				m.environmentCursor = 0
				m.loading = false
				m.err = nil

			case resourceDetailsScreen:
				m.screen = resourcesScreen
				m.loading = false
				m.err = nil

			case deploymentsScreen:
				m.screen = resourceDetailsScreen
				m.deployments = nil
				m.deploymentCount = 0
				m.deploymentCursor = 0
				m.deploymentSkip = 0
				m.loading = false
				m.err = nil

			case deploymentDetailsScreen:
				m.screen = deploymentsScreen
				m.loading = false
				m.err = nil
			}

		case "r":
			m.loading = true
			m.err = nil

			if m.screen == projectsScreen {
				return m, m.loadProjects()
			}

			if m.screen == resourcesScreen &&
				m.project != nil &&
				len(m.project.Environments) > 0 {
				environment := m.project.Environments[m.environmentCursor]
				return m, m.loadResources(environment.ID)
			}

			if m.screen == deploymentsScreen {
				resource := m.selectedResource()

				if resource != nil {
					return m, m.loadDeployments(
						resource.UUID,
						m.deploymentSkip,
					)
				}
			}

			if m.project != nil {
				return m, m.loadProject(m.project.UUID)
			}

		case "enter":
			if m.loading || m.err != nil {
				break
			}

			switch m.screen {
			case projectsScreen:
				if len(m.projects) == 0 {
					break
				}

				project := m.projects[m.projectCursor]
				m.loading = true
				m.err = nil

				return m, m.loadProject(project.UUID)

			case environmentsScreen:
				if m.project == nil ||
					len(m.project.Environments) == 0 {
					break
				}

				environment := m.project.Environments[m.environmentCursor]
				m.loading = true
				m.err = nil

				return m, m.loadResources(environment.ID)

			case resourcesScreen:
				if m.selectedResource() != nil {
					m.screen = resourceDetailsScreen
				}

			case deploymentsScreen:
				if m.selectedDeployment() != nil {
					m.screen = deploymentDetailsScreen
				}
			}

		case "up", "k":
			m.moveCursor(-1)

		case "down", "j":
			m.moveCursor(1)

		case "g":
			m.moveToFirst()

		case "G":
			m.moveToLast()

		case "d":
			if m.screen != resourceDetailsScreen ||
				m.loading ||
				m.err != nil {
				break
			}

			resource := m.selectedResource()
			if resource == nil ||
				!strings.EqualFold(resource.Type, "application") {
				break
			}

			m.loading = true
			m.err = nil
			m.deploymentSkip = 0
			return m, m.loadDeployments(resource.UUID, 0)

		case "n":
			if m.screen != deploymentsScreen ||
				m.loading {
				break
			}
			nextSkip := m.deploymentSkip + m.deploymentTake
			if nextSkip >= m.deploymentCount {
				break
			}
			resource := m.selectedResource()
			if resource == nil {
				break
			}
			m.loading = true
			m.err = nil
			return m, m.loadDeployments(
				resource.UUID,
				nextSkip,
			)

		case "p":
			if m.screen != deploymentsScreen ||
				m.loading {
				break
			}
			prevSkip := m.deploymentSkip - m.deploymentTake
			if prevSkip < 0 {
				break
			}
			resource := m.selectedResource()
			if resource == nil {
				break
			}
			m.loading = true
			m.err = nil
			return m, m.loadDeployments(
				resource.UUID,
				prevSkip,
			)
		}
	}

	return m, nil
}

func (m *Model) moveCursor(change int) {
	switch m.screen {
	case projectsScreen:
		next := m.projectCursor + change

		if next >= 0 && next < len(m.projects) {
			m.projectCursor = next
		}

	case environmentsScreen:
		if m.project == nil {
			return
		}

		next := m.environmentCursor + change

		if next >= 0 &&
			next < len(m.project.Environments) {
			m.environmentCursor = next
		}

	case resourcesScreen:
		next := m.resourceCursor + change

		if next >= 0 && next < len(m.resources) {
			m.resourceCursor = next
		}

	case deploymentsScreen:
		next := m.deploymentCursor + change

		if next >= 0 && next < len(m.deployments) {
			m.deploymentCursor = next
		}
	}
}

func (m *Model) moveToFirst() {
	switch m.screen {
	case projectsScreen:
		m.projectCursor = 0

	case environmentsScreen:
		m.environmentCursor = 0

	case resourcesScreen:
		m.resourceCursor = 0

	case deploymentsScreen:
		m.deploymentCursor = 0
	}
}

func (m *Model) moveToLast() {
	switch m.screen {
	case projectsScreen:
		if len(m.projects) > 0 {
			m.projectCursor = len(m.projects) - 1
		}

	case environmentsScreen:
		if m.project != nil &&
			len(m.project.Environments) > 0 {
			m.environmentCursor =
				len(m.project.Environments) - 1
		}

	case resourcesScreen:
		if len(m.resources) > 0 {
			m.resourceCursor = len(m.resources) - 1
		}

	case deploymentsScreen:
		if len(m.deployments) > 0 {
			m.deploymentCursor = len(m.deployments) - 1
		}
	}
}
