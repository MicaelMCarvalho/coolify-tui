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

type APIError struct {
	Message string `json:"message"`
}
