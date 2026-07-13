package ui

import tea "github.com/charmbracelet/bubbletea"

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
			}

		case "up", "k":
			m.moveCursor(-1)

		case "down", "j":
			m.moveCursor(1)

		case "g":
			m.moveToFirst()

		case "G":
			m.moveToLast()
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
	}
}
