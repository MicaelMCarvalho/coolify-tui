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

	case errMsg:
		m.loading = false
		m.err = msg.err

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "esc":
			if m.screen == environmentsScreen {
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

			if m.project != nil {
				return m, m.loadProject(m.project.UUID)
			}

		case "enter":
			if m.loading || m.err != nil {
				break
			}

			if m.screen == projectsScreen &&
				len(m.projects) > 0 {
				project := m.projects[m.projectCursor]

				m.loading = true
				m.err = nil

				return m, m.loadProject(project.UUID)
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
	}
}

func (m *Model) moveToFirst() {
	switch m.screen {
	case projectsScreen:
		m.projectCursor = 0

	case environmentsScreen:
		m.environmentCursor = 0
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
	}
}
