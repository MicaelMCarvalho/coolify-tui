package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/micaelmcarvalho/coolify-tui/internal/coolify"
)

type panel int

const (
	teamsPanel panel = iota
	projectsPanel
	environmentsPanel
	resourcesPanel
	detailsPanel
	environmentVariablesPanel
	deploymentsPanel
)

type teamsLoadedMsg struct {
	teams []coolify.Team
}

type projectsLoadedMsg struct {
	projects []coolify.Project
}

type projectLoadedMsg struct {
	projectUUID string
	project     coolify.ProjectDetails
}

type resourcesLoadedMsg struct {
	environmentID int
	resources     []coolify.Resource
}

type environmentVariablesLoadedMsg struct {
	resourceUUID string
	variables    []coolify.EnvironmentVariable
}

type deploymentsLoadedMsg struct {
	resourceUUID string
	result       coolify.DeploymentList
	skip         int
}

type deploymentDetailsLoadedMsg struct {
	resourceUUID string
	deployment   coolify.DeploymentDetails
}

type deploymentStartedMsg struct {
	resourceUUID string
	result       coolify.StartDeploymentResult
}

type deploymentDetailsFailedMsg struct {
	deploymentUUID string
	err            error
}

type deploymentPollMsg struct {
	deploymentUUID string
}

type errMsg struct {
	err error
}

type Model struct {
	// client      *coolify.Client
	client      APIClient
	activePanel panel

	teams      []coolify.Team
	teamCursor int

	projects      []coolify.Project
	projectCursor int

	project           *coolify.ProjectDetails
	environmentCursor int

	resources      []coolify.Resource
	resourceCursor int

	deployments      []coolify.Deployment
	deploymentCount  int
	deploymentCursor int
	deploymentSkip   int
	deploymentTake   int

	deploymentDetailsOpen bool
	deploymentDetailsUUID string
	deploymentDetails     *coolify.DeploymentDetails
	deploymentLogOffset   int
	deploymentFollowLogs  bool
	deploymentPolling     bool
	deploymentPollPending bool
	deployConfirmOpen     bool

	environmentVariables       []coolify.EnvironmentVariable
	environmentVariablesCursor int
	revealEnvironmentValues    bool

	filtering      bool
	filterPanel    panel
	filterInput    string
	filterOriginal string
	filters        map[panel]string

	helpOpen bool

	width   int
	height  int
	loading bool
	err     error
}

func NewModel(client APIClient) Model {
	return Model{
		client:         client,
		activePanel:    teamsPanel,
		loading:        true,
		deploymentTake: 20,
		teams:          []coolify.Team{},
		projects:       []coolify.Project{},
		resources:      []coolify.Resource{},
		deployments:    []coolify.Deployment{},
		filters:        make(map[panel]string),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadTeams(),
		m.loadProjects(),
	)
}

func (m Model) loadTeams() tea.Cmd {
	return func() tea.Msg {
		teams, err := m.client.ListTeams(
			context.Background(),
		)
		if err != nil {
			return errMsg{err: err}
		}

		return teamsLoadedMsg{teams: teams}
	}
}

func (m Model) loadProjects() tea.Cmd {
	return func() tea.Msg {
		projects, err := m.client.ListProjects(
			context.Background(),
		)
		if err != nil {
			return errMsg{err: err}
		}

		return projectsLoadedMsg{
			projects: projects,
		}
	}
}

func (m Model) loadProject(uuid string) tea.Cmd {
	return func() tea.Msg {
		project, err := m.client.GetProject(
			context.Background(),
			uuid,
		)
		if err != nil {
			return errMsg{err: err}
		}

		return projectLoadedMsg{
			projectUUID: uuid,
			project:     project,
		}
	}
}

func (m Model) selectedTeam() *coolify.Team {
	if len(m.teams) == 0 ||
		m.teamCursor < 0 ||
		m.teamCursor >= len(m.teams) {
		return nil
	}

	return &m.teams[m.teamCursor]
}

func (m Model) selectedProject() *coolify.Project {
	if len(m.projects) == 0 ||
		m.projectCursor < 0 ||
		m.projectCursor >= len(m.projects) {
		return nil
	}

	return &m.projects[m.projectCursor]
}

func (m Model) selectedEnvironment() *coolify.Environment {
	if m.project == nil ||
		len(m.project.Environments) == 0 ||
		m.environmentCursor < 0 ||
		m.environmentCursor >= len(m.project.Environments) {
		return nil
	}

	return &m.project.Environments[m.environmentCursor]
}

func (m Model) selectedResource() *coolify.Resource {
	if len(m.resources) == 0 ||
		m.resourceCursor < 0 ||
		m.resourceCursor >= len(m.resources) {
		return nil
	}

	return &m.resources[m.resourceCursor]
}

func (m Model) loadResources(
	environmentID int,
) tea.Cmd {
	return func() tea.Msg {
		resources, err := m.client.ListResources(
			context.Background(),
			environmentID,
		)
		if err != nil {
			return errMsg{err: err}
		}

		return resourcesLoadedMsg{
			environmentID: environmentID,
			resources:     resources,
		}
	}
}

func (m Model) loadDeployments(
	applicationUUID string,
	skip int,
) tea.Cmd {
	return func() tea.Msg {
		result, err := m.client.ListApplicationDeployments(
			context.Background(),
			applicationUUID,
			skip,
			m.deploymentTake,
		)
		if err != nil {
			return errMsg{err: err}
		}
		return deploymentsLoadedMsg{
			resourceUUID: applicationUUID,
			result:       result,
			skip:         skip,
		}
	}
}

func (m Model) loadDeploymentDetails(
	deploymentUUID string,
) tea.Cmd {
	return func() tea.Msg {
		details, err := m.client.GetDeployment(
			context.Background(),
			deploymentUUID,
		)
		if err != nil {
			return deploymentDetailsFailedMsg{
				deploymentUUID: deploymentUUID,
				err:            err,
			}
		}

		return deploymentDetailsLoadedMsg{
			resourceUUID: deploymentUUID,
			deployment:   details,
		}
	}
}

func (m Model) startApplicationDeployment(
	applicationUUID string,
) tea.Cmd {
	return func() tea.Msg {
		result, err := m.client.StartApplicationDeployment(
			context.Background(),
			applicationUUID,
			false,
		)
		if err != nil {
			return errMsg{err: err}
		}

		return deploymentStartedMsg{
			resourceUUID: applicationUUID,
			result:       result,
		}
	}
}

func (m Model) loadEnvironmentVariables(
	applicationUUID string,
) tea.Cmd {
	return func() tea.Msg {
		variables, err :=
			m.client.ListApplicationEnvironmentVariables(
				context.Background(),
				applicationUUID,
			)
		if err != nil {
			return errMsg{err: err}
		}

		return environmentVariablesLoadedMsg{
			resourceUUID: applicationUUID,
			variables:    variables,
		}
	}
}

func (m Model) selectedDeployment() *coolify.Deployment {
	if len(m.deployments) == 0 ||
		m.deploymentCursor < 0 ||
		m.deploymentCursor >= len(m.deployments) {
		return nil
	}

	return &m.deployments[m.deploymentCursor]
}

const deploymentPollInterval = 2 * time.Second

func pollDeploymentAfter(
	deploymentUUID string,
) tea.Cmd {
	return tea.Tick(
		deploymentPollInterval,
		func(time.Time) tea.Msg {
			return deploymentPollMsg{
				deploymentUUID: deploymentUUID,
			}
		},
	)
}
