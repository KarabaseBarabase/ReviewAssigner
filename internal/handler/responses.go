package handler

import "ReviewAssigner/internal/models"

type PRResponse struct {
	PR *models.PullRequest `json:"pr"`
}

type TeamResponse struct {
	Team *models.Team `json:"team"`
}

type UserResponse struct {
	User *models.User `json:"user"`
}

type UserPRsResponse struct {
	UserID       string                    `json:"user_id"`
	PullRequests []models.PullRequestShort `json:"pull_requests"`
}

type DeactivateUsersResponse struct {
	TeamName string                 `json:"team_name"`
	Results  map[string]interface{} `json:"results"`
}

type StatsResponse struct {
	UserAssignments map[string]int         `json:"user_assignments,omitempty"`
	PRMetrics       map[string]interface{} `json:"pr_metrics,omitempty"`
}

type ReassignReviewerResponse struct {
	PR         *models.PullRequest `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}
