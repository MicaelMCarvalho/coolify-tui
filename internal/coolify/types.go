package coolify

import "time"

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

type APIError struct {
	Message string `json:"message"`
}
