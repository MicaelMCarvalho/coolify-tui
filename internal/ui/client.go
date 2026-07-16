package ui

import (
	"context"

	"github.com/micaelmcarvalho/coolify-tui/internal/coolify"
)

type APIClient interface {
	ListTeams(context.Context) ([]coolify.Team, error)
	ListProjects(context.Context) ([]coolify.Project, error)

	GetProject(
		context.Context,
		string,
	) (coolify.ProjectDetails, error)

	ListResources(
		context.Context,
		int,
	) ([]coolify.Resource, error)

	ListApplicationEnvironmentVariables(
		context.Context,
		string,
	) ([]coolify.EnvironmentVariable, error)

	ListApplicationDeployments(
		context.Context,
		string,
		int,
		int,
	) (coolify.DeploymentList, error)

	GetDeployment(
		context.Context,
		string,
	) (coolify.DeploymentDetails, error)

	StartApplicationDeployment(
		context.Context,
		string,
		bool,
	) (coolify.StartDeploymentResult, error)
}
