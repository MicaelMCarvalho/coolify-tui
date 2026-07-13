package coolify

import (
	"encoding/json"
	"testing"
)

func TestResourceEnvironmentID(t *testing.T) {
	resource := Resource{
		ID:            1,
		Name:          "example",
		EnvironmentID: 42,
	}

	if resource.EnvironmentID != 42 {
		t.Fatalf(
			"Expected EnvironmentID to be 42, got %d",
			resource.EnvironmentID,
		)
	}
}

func TestDeploymentListDecoding(t *testing.T) {
	input := []byte(`{
		"count": 1,
		"deployments": [
			{
				"id": 10,
				"deployment_uuid": "deployment-123",
				"application_id": 5,
				"application_name": "example",
				"status": "finished",
				"commit": "abc123",
				"commit_message": "Deploy example",
				"server_name": "production",
				"created_at": "2026-07-13T10:00:00Z",
				"finished_at": "2026-07-13T10:01:00Z"
			}
		]
	}`)

	var result DeploymentList

	if err := json.Unmarshal(input, &result); err != nil {
		t.Fatalf("decode deployment list: %v", err)
	}

	if result.Count != 1 {
		t.Fatalf("expected count 1, got %d", result.Count)
	}

	if len(result.Deployments) != 1 {
		t.Fatalf(
			"expected one deployment, got %d",
			len(result.Deployments),
		)
	}

	if result.Deployments[0].Status != "finished" {
		t.Fatalf(
			"expected finished status, got %q",
			result.Deployments[0].Status,
		)
	}
}
