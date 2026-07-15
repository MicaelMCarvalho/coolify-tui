package coolify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	baseURL = strings.TrimSpace(baseURL)
	baseURL = strings.TrimRight(baseURL, "/")
	baseURL = strings.TrimSuffix(baseURL, "/api/v1")

	return &Client{
		baseURL: baseURL + "/api/v1",
		token:   strings.TrimSpace(token),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) ListTeams(ctx context.Context) ([]Team, error) {
	var raw json.RawMessage

	if err := c.get(ctx, "/teams", &raw); err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}

	// Newer/other Coolify versions may return a normal JSON array.
	var teams []Team
	if err := json.Unmarshal(raw, &teams); err == nil {
		return teams, nil
	}

	var teamMap map[string]Team
	if err := json.Unmarshal(raw, &teamMap); err != nil {
		return nil, fmt.Errorf("decode teams response: %w", err)
	}

	teams = make([]Team, 0, len(teamMap))
	for _, team := range teamMap {
		teams = append(teams, team)
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})

	return teams, nil
}

func (c *Client) ListProjects(ctx context.Context) ([]Project, error) {
	var projects []Project

	if err := c.get(ctx, "/projects", &projects); err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	sort.Slice(projects, func(i, j int) bool {
		return strings.ToLower(projects[i].Name) <
			strings.ToLower(projects[j].Name)
	})

	return projects, nil
}

func (c *Client) GetProject(
	ctx context.Context,
	projectUUID string,
) (ProjectDetails, error) {
	var project ProjectDetails

	projectUUID = strings.TrimSpace(projectUUID)
	if projectUUID == "" {
		return ProjectDetails{}, fmt.Errorf("project UUID is empty")
	}

	path := "/projects/" + url.PathEscape(projectUUID)

	if err := c.get(ctx, path, &project); err != nil {
		return ProjectDetails{}, fmt.Errorf("get project: %w", err)
	}

	return project, nil
}

func (c *Client) ListResources(
	ctx context.Context,
	environmentID int,
) ([]Resource, error) {
	if environmentID <= 0 {
		return nil, fmt.Errorf("environment ID must be positive")
	}

	var allResources []Resource

	if err := c.get(ctx, "/resources", &allResources); err != nil {
		return nil, fmt.Errorf("list resources: %w", err)
	}

	resources := make([]Resource, 0)

	for _, resource := range allResources {
		if resource.EnvironmentID == environmentID {
			resources = append(resources, resource)
		}
	}

	sort.Slice(resources, func(i, j int) bool {
		left := strings.ToLower(resources[i].Name)
		right := strings.ToLower(resources[j].Name)

		if left == right {
			return resources[i].Type < resources[j].Type
		}
		return left < right
	})

	return resources, nil
}

func (c *Client) ListApplicationDeployments(
	ctx context.Context,
	applicationUUID string,
	skip int,
	take int,
) (DeploymentList, error) {
	applicationUUID = strings.TrimSpace(applicationUUID)

	if applicationUUID == "" {
		return DeploymentList{},
			fmt.Errorf("application UUID is required")
	}

	if skip < 0 {
		return DeploymentList{},
			fmt.Errorf("skip cannot be negative")
	}

	if take < 1 {
		return DeploymentList{},
			fmt.Errorf("take must be positive")
	}

	query := url.Values{}
	query.Set("skip", strconv.Itoa(skip))
	query.Set("take", strconv.Itoa(take))

	path := "/deployments/applications/" +
		url.PathEscape(applicationUUID) +
		"?" +
		query.Encode()

	var deployments DeploymentList

	if err := c.get(ctx, path, &deployments); err != nil {
		return DeploymentList{},
			fmt.Errorf("list application deployments: %w", err)
	}

	return deployments, nil
}

func (c *Client) GetDeployment(
	ctx context.Context,
	deploymentUUID string,
) (DeploymentDetails, error) {
	deploymentUUID = strings.TrimSpace(deploymentUUID)

	if deploymentUUID == "" {
		return DeploymentDetails{},
			fmt.Errorf("deployment UUID is required")
	}

	path := "/deployments/" +
		url.PathEscape(deploymentUUID)

	var deployment DeploymentDetails

	if err := c.get(ctx, path, &deployment); err != nil {
		return DeploymentDetails{},
			fmt.Errorf("get deployment: %w", err)
	}

	return deployment, nil
}

func (c *Client) ListApplicationEnvironmentVariables(
	ctx context.Context,
	applicationUUID string,
) ([]EnvironmentVariable, error) {
	applicationUUID = strings.TrimSpace(applicationUUID)

	if applicationUUID == "" {
		return nil, fmt.Errorf("application UUID is required")
	}

	path := "/applications/" +
		url.PathEscape(applicationUUID) +
		"/envs"

	var variables []EnvironmentVariable

	if err := c.get(ctx, path, &variables); err != nil {
		return nil, fmt.Errorf(
			"list application environment variables: %w",
			err,
		)
	}

	sort.SliceStable(variables, func(i, j int) bool {
		return strings.ToLower(variables[i].Key) <
			strings.ToLower(variables[j].Key)
	})

	return variables, nil
}

func (c *Client) get(
	ctx context.Context,
	path string,
	result any,
) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+path,
		nil,
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "coolify-tui")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return decodeAPIError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func decodeAPIError(resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return fmt.Errorf("Coolify returned HTTP %d", resp.StatusCode)
	}

	var apiErr APIError

	if json.Unmarshal(body, &apiErr) == nil && apiErr.Message != "" {
		return fmt.Errorf(
			"Coolify returned HTTP %d: %s",
			resp.StatusCode,
			apiErr.Message,
		)
	}

	return fmt.Errorf(
		"Coolify returned HTTP %d: %s",
		resp.StatusCode,
		strings.TrimSpace(string(body)),
	)
}
