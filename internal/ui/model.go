package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/micaelmcarvalho/coolify-tui/internal/coolify"
)

type projectsLoadedMsg struct {
	projects []coolify.Project
}

type errMsg struct {
	err error
}

type Model struct {
	client   *coolify.Client
	projects []coolify.Project
	cursor   int
	width    int
	height   int
	loading  bool
	err      error
}

func NewModel(client *coolify.Client) Model {
	return Model{
		client:  client,
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadProjects()
}

func (m Model) loadProjects() tea.Cmd {
	return func() tea.Msg {
		projects, err := m.client.ListProjects(context.Background())
		if err != nil {
			return errMsg{err}
		}
		return projectsLoadedMsg{projects}
	}
}
