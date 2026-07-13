package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/micaelmcarvalho/coolify-tui/internal/coolify"
)

type screen int

const (
	projectsScreen screen = iota
	environmentsScreen
	resourcesScreen
)

type projectsLoadedMsg struct {
	projects []coolify.Project
}

type projectLoadedMsg struct {
	project coolify.ProjectDetails
}

type resourcesLoadedMsg struct {
	resources []coolify.Resource
}

type errMsg struct {
	err error
}

type Model struct {
	client *coolify.Client
	screen screen

	projects      []coolify.Project
	projectCursor int

	project           *coolify.ProjectDetails
	environmentCursor int

	resources      []coolify.Resource
	resourceCursor int

	width   int
	height  int
	loading bool
	err     error
}

func NewModel(client *coolify.Client) Model {
	return Model{
		client:  client,
		screen:  projectsScreen,
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

func (m Model) loadProject(uuid string) tea.Cmd {
	return func() tea.Msg {
		project, err := m.client.GetProject(
			context.Background(),
			uuid,
		)
		if err != nil {
			return errMsg{err}
		}
		return projectLoadedMsg{project: project}
	}
}

func (m Model) loadResources(environmentID int) tea.Cmd {
	return func() tea.Msg {
		resources, err := m.client.ListResources(
			context.Background(),
			environmentID,
		)
		if err != nil {
			return errMsg{err: err}
		}
		return resourcesLoadedMsg{resources: resources}
	}
}
