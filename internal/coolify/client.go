package coolify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
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
