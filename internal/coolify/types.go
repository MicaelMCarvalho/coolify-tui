package coolify

import (
	"encoding/json"
	"time"
)

type Team struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Description  *string    `json:"description"`
	PersonalTeam bool       `json:"personal_team"`
	ShowBoarding bool       `json:"show_boarding"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type Project struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Environment struct {
	ID          int        `json:"id"`
	UUID        string     `json:"uuid"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	ProjectID   int        `json:"project_id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ProjectDetails struct {
	ID           int           `json:"id"`
	UUID         string        `json:"uuid"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Environments []Environment `json:"environments"`
	CreatedAt    *time.Time    `json:"created_at"`
	UpdatedAt    *time.Time    `json:"updated_at"`
}

type Resource struct {
	ID            int     `json:"id"`
	UUID          string  `json:"uuid"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	Status        string  `json:"status"`
	EnvironmentID int     `json:"environment_id"`
	Description   *string `json:"description"`
	FQDN          *string `json:"fqdn"`
}

type Deployment struct {
	// ID              int     `json:"id"`
	DeploymentUUID string `json:"deployment_uuid"`
	// ApplicationID   int     `json:"application_id"`
	ApplicationName string  `json:"application_name"`
	Status          string  `json:"status"`
	Commit          string  `json:"commit"`
	CommitMessage   *string `json:"commit_message"`
	ServerName      string  `json:"server_name"`
	DeploymentURL   *string `json:"deployment_url"`
	ForceRebuild    bool    `json:"force_rebuild"`
	RestartOnly     bool    `json:"restart_only"`
	Rollback        bool    `json:"rollback"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	FinishedAt      *string `json:"finished_at"`
}

type DeploymentList struct {
	Count       int          `json:"count"`
	Deployments []Deployment `json:"deployments"`
}

type DeploymentDetails struct {
	DeploymentUUID  string          `json:"deployment_uuid"`
	ApplicationName string          `json:"application_name"`
	Status          string          `json:"status"`
	Commit          string          `json:"commit"`
	CommitMessage   *string         `json:"commit_message"`
	ServerName      string          `json:"server_name"`
	DeploymentURL   *string         `json:"deployment_url"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
	FinishedAt      *string         `json:"finished_at"`
	Logs            json.RawMessage `json:"logs"`
}

type EnvironmentVariable struct {
	UUID        string  `json:"uuid"`
	Key         string  `json:"key"`
	Value       string  `json:"value"`
	RealValue   string  `json:"real_value"`
	Comment     *string `json:"comment"`
	IsBuildTime bool    `json:"is_buildtime"`
	IsRuntime   bool    `json:"is_runtime"`
	IsPreview   bool    `json:"is_preview"`
	IsLiteral   bool    `json:"is_literal"`
	IsMultiline bool    `json:"is_multiline"`
	IsShared    bool    `json:"is_shared"`
	IsShownOnce bool    `json:"is_shown_once"`
	IsCoolify   bool    `json:"is_coolify"`
}

type APIError struct {
	Message string `json:"message"`
}
