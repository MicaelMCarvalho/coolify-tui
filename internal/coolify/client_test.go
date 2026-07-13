package coolify

import "testing"

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
