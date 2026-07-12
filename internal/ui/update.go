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

		if m.cursor >= len(m.projects) {
			m.cursor = max(0, len(m.projects)-1)
		}

	case errMsg:
		m.err = msg.err
		m.loading = false

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.projects)-1 {
				m.cursor++
			}
		case "g":
			m.cursor = 0
		case "G":
			if len(m.projects) > 0 {
				m.cursor = len(m.projects) - 1
			}
		case "r":
			m.loading = true
			m.err = nil
			return m, m.loadProjects()
		}
	}

	return m, nil
}
