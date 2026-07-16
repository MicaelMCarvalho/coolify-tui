package coolify

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

type DemoClient struct {
	mutex     sync.Mutex
	pollCount int
}

func NewDemoClient() *DemoClient {
	return &DemoClient{}
}

func (c *DemoClient) ListTeams(
	context.Context,
) ([]Team, error) {
	return []Team{
		{
			ID:           1,
			Name:         "Demo Platform Team",
			PersonalTeam: true,
		},
		{
			ID:   2,
			Name: "Operations",
		},
	}, nil
}

func (c *DemoClient) ListProjects(
	context.Context,
) ([]Project, error) {
	return demoProjects(), nil
}

func (c *DemoClient) GetProject(
	_ context.Context,
	projectUUID string,
) (ProjectDetails, error) {
	for _, project := range demoProjects() {
		if project.UUID != projectUUID {
			continue
		}

		return ProjectDetails{
			ID:          project.ID,
			UUID:        project.UUID,
			Name:        project.Name,
			Description: project.Description,
			Environments: []Environment{
				{
					ID:        project.ID*100 + 1,
					UUID:      project.UUID + "-production",
					Name:      "production",
					ProjectID: project.ID,
				},
				{
					ID:        project.ID*100 + 2,
					UUID:      project.UUID + "-staging",
					Name:      "staging",
					ProjectID: project.ID,
				},
			},
		}, nil
	}

	return ProjectDetails{},
		fmt.Errorf("demo project not found")
}

func (c *DemoClient) ListResources(
	_ context.Context,
	environmentID int,
) ([]Resource, error) {
	suffix := fmt.Sprintf("%d", environmentID)

	apiDescription := "Public application API"
	workerDescription := "Background job processor"
	databaseDescription := "Application database"
	redisDescription := "Cache and queue service"
	apiURL := "https://api.demo.example.com"

	return []Resource{
		{
			ID:            environmentID*10 + 1,
			UUID:          "demo-api-" + suffix,
			Name:          "demo-api",
			Type:          "application",
			Status:        "running:healthy",
			GitBranch:     "main",
			EnvironmentID: environmentID,
			Description:   &apiDescription,
			FQDN:          &apiURL,
		},
		{
			ID:            environmentID*10 + 2,
			UUID:          "demo-worker-" + suffix,
			Name:          "background-worker",
			Type:          "application",
			Status:        "running:healthy",
			GitBranch:     "main",
			EnvironmentID: environmentID,
			Description:   &workerDescription,
		},
		{
			ID:            environmentID*10 + 3,
			UUID:          "demo-postgres-" + suffix,
			Name:          "postgres",
			Type:          "database",
			Status:        "running:healthy",
			EnvironmentID: environmentID,
			Description:   &databaseDescription,
		},
		{
			ID:            environmentID*10 + 4,
			UUID:          "demo-redis-" + suffix,
			Name:          "redis",
			Type:          "service",
			Status:        "running:healthy",
			EnvironmentID: environmentID,
			Description:   &redisDescription,
		},
	}, nil
}

func (c *DemoClient) ListApplicationEnvironmentVariables(
	context.Context,
	string,
) ([]EnvironmentVariable, error) {
	return []EnvironmentVariable{
		demoVariable("APP_ENV", "production", true, true),
		demoVariable("APP_URL", "https://demo.example.com", true, true),
		demoVariable("DATABASE_HOST", "postgres", true, true),
		demoVariable("DATABASE_NAME", "demo", true, true),
		demoVariable("DATABASE_USER", "demo_user", true, true),
		demoVariable("DATABASE_PASSWORD", "fake-demo-password", true, true),
		demoVariable("REDIS_HOST", "redis", true, true),
		demoVariable("LOG_LEVEL", "info", true, true),
		demoVariable("FEATURE_REPORTS", "true", true, true),
		demoVariable("API_TOKEN", "fake-demo-token", false, true),
	}, nil
}

func (c *DemoClient) ListApplicationDeployments(
	_ context.Context,
	applicationUUID string,
	skip int,
	take int,
) (DeploymentList, error) {
	finishedOne := "2026-07-16T14:34:22Z"
	finishedTwo := "2026-07-15T18:12:04Z"
	finishedThree := "2026-07-14T09:45:11Z"

	deployments := []Deployment{
		{
			DeploymentUUID:  applicationUUID + "-deploy-1",
			ApplicationName: "demo-api",
			Status:          "finished",
			Commit:          "c4f31f28a9d89eab",
			CommitMessage:   demoString("Add live deployment tracking"),
			ServerName:      "demo-server-01",
			CreatedAt:       "2026-07-16T14:32:02Z",
			UpdatedAt:       finishedOne,
			FinishedAt:      &finishedOne,
		},
		{
			DeploymentUUID:  applicationUUID + "-deploy-2",
			ApplicationName: "demo-api",
			Status:          "finished",
			Commit:          "96d80c2a31ef725a",
			CommitMessage:   demoString("Improve dashboard navigation"),
			ServerName:      "demo-server-01",
			CreatedAt:       "2026-07-15T18:09:20Z",
			UpdatedAt:       finishedTwo,
			FinishedAt:      &finishedTwo,
		},
		{
			DeploymentUUID:  applicationUUID + "-deploy-3",
			ApplicationName: "demo-api",
			Status:          "failed",
			Commit:          "448c01f7f6223ce1",
			CommitMessage:   demoString("Update application dependencies"),
			ServerName:      "demo-server-01",
			CreatedAt:       "2026-07-14T09:42:40Z",
			UpdatedAt:       finishedThree,
			FinishedAt:      &finishedThree,
		},
	}

	count := len(deployments)

	if skip >= count {
		return DeploymentList{
			Count:       count,
			Deployments: []Deployment{},
		}, nil
	}

	end := min(skip+take, count)

	return DeploymentList{
		Count:       count,
		Deployments: deployments[skip:end],
	}, nil
}

func (c *DemoClient) GetDeployment(
	_ context.Context,
	deploymentUUID string,
) (DeploymentDetails, error) {
	status := "finished"
	finishedAt := demoString(
		"2026-07-16T14:34:22Z",
	)

	logMessages := []string{
		"Deployment request accepted.",
		"Cloning repository.",
		"Building application image.",
		"Starting new container.",
		"Health check passed.",
		"Deployment completed successfully.",
	}

	if deploymentUUID == "demo-live-deployment" {
		c.mutex.Lock()
		c.pollCount++
		pollCount := c.pollCount
		c.mutex.Unlock()

		switch {
		case pollCount == 1:
			status = "queued"
			finishedAt = nil
			logMessages = logMessages[:1]

		case pollCount < 4:
			status = "in_progress"
			finishedAt = nil
			logMessages = logMessages[:pollCount+1]

		default:
			status = "finished"
		}
	}

	return DeploymentDetails{
		DeploymentUUID:  deploymentUUID,
		ApplicationName: "demo-api",
		Status:          status,
		Commit:          "c4f31f28a9d89eab",
		CommitMessage: demoString(
			"Add live deployment tracking",
		),
		ServerName: "demo-server-01",
		CreatedAt:  "2026-07-16T14:32:02Z",
		UpdatedAt:  "2026-07-16T14:34:22Z",
		FinishedAt: finishedAt,
		Logs:       demoLogs(logMessages),
	}, nil
}

func (c *DemoClient) StartApplicationDeployment(
	context.Context,
	string,
	bool,
) (StartDeploymentResult, error) {
	c.mutex.Lock()
	c.pollCount = 0
	c.mutex.Unlock()

	return StartDeploymentResult{
		Message:        "Demo deployment queued.",
		DeploymentUUID: "demo-live-deployment",
	}, nil
}

func demoProjects() []Project {
	return []Project{
		{
			ID:          1,
			UUID:        "demo-atlas",
			Name:        "Atlas API",
			Description: "Public API and processing services",
		},
		{
			ID:          2,
			UUID:        "demo-nova",
			Name:        "Nova Dashboard",
			Description: "Operations monitoring dashboard",
		},
		{
			ID:          3,
			UUID:        "demo-storefront",
			Name:        "Storefront",
			Description: "Customer-facing web application",
		},
		{
			ID:          4,
			UUID:        "demo-tools",
			Name:        "Internal Tools",
			Description: "Internal automation utilities",
		},
	}
}

func demoVariable(
	key string,
	value string,
	buildTime bool,
	runtime bool,
) EnvironmentVariable {
	return EnvironmentVariable{
		UUID:        "demo-" + key,
		Key:         key,
		Value:       value,
		RealValue:   value,
		IsBuildTime: buildTime,
		IsRuntime:   runtime,
	}
}

func demoLogs(
	messages []string,
) json.RawMessage {
	type logEntry struct {
		Output    string `json:"output"`
		Timestamp string `json:"timestamp"`
		Hidden    bool   `json:"hidden"`
	}

	entries := make(
		[]logEntry,
		0,
		len(messages),
	)

	for index, message := range messages {
		entries = append(entries, logEntry{
			Output: message,
			Timestamp: fmt.Sprintf(
				"2026-07-16T14:32:%02dZ",
				index*5,
			),
		})
	}

	data, _ := json.Marshal(entries)

	return data
}

func demoString(value string) *string {
	return &value
}
