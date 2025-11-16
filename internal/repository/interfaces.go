package repository

import "ReviewAssigner/internal/models"

type UserRepository interface {
	CreateOrUpdateUser(user *models.User) error
	GetUserByID(userID string) (*models.User, error)
	SetUserActive(userID string, isActive bool) error
	GetActiveTeamMembers(teamName string, excludeUserID string) ([]models.User, error)
	GetUsersByTeam(teamName string) ([]models.User, error)
}

type TeamRepository interface {
	CreateTeam(teamName string) error
	TeamExists(teamName string) (bool, error)
	GetTeam(teamName string) (*models.Team, error)
	GetUsersByTeam(teamName string) ([]models.User, error)
}

type PRRepository interface {
	CreatePR(pr *models.PullRequest) error
	PRExists(prID string) (bool, error)
	GetPRByID(prID string) (*models.PullRequest, error)
	MergePR(prID string) error
	AddPRReviewer(prID, reviewerID string) error
	ReplacePRReviewer(prID, oldReviewerID, newReviewerID string) error
	GetPRReviewers(prID string) ([]string, error)
	IsReviewerAssigned(prID, reviewerID string) (bool, error)
	GetAssignedPRs(userID string) ([]models.PullRequestShort, error)
	GetUserAssignmentStats() (map[string]int, error)
	GetPRMetrics() (map[string]interface{}, error)
	DeletePR(prID string) error //new
}

type ReviewService interface {
	AssignReviewers(teamName, authorID, prID string) ([]string, error)
	ReplaceReviewer(prID, oldReviewerID string) (string, error)
}
