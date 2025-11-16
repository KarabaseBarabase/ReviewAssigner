package repository

import (
	"fmt"

	"ReviewAssigner/internal/models"

	"github.com/jmoiron/sqlx"
)

// реализует TeamRepository интерфейс
type TeamRepositoryImpl struct {
	db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) *TeamRepositoryImpl {
	return &TeamRepositoryImpl{db: db}
}

func (r *TeamRepositoryImpl) CreateTeam(teamName string) error {
	query := `INSERT INTO teams (team_name) VALUES ($1)`
	_, err := r.db.Exec(query, teamName)
	return err
}

func (r *TeamRepositoryImpl) TeamExists(teamName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`
	err := r.db.Get(&exists, query, teamName)
	return exists, err
}

func (r *TeamRepositoryImpl) GetTeam(teamName string) (*models.Team, error) {
	exists, err := r.TeamExists(teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("team '%s' not found", teamName)
	}

	users, err := r.GetUsersByTeam(teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team users: %w", err)
	}

	members := make([]models.TeamMember, len(users))
	for i, user := range users {
		members[i] = models.TeamMember{
			UserID:   user.UserID,
			Username: user.Username,
			IsActive: user.IsActive,
		}
	}

	return &models.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}

func (r *TeamRepositoryImpl) GetUsersByTeam(teamName string) ([]models.User, error) {
	var users []models.User
	query := `
        SELECT 
            user_id, 
            username, 
            team_name, 
            is_active, 
            created_at, 
            updated_at
        FROM users 
        WHERE team_name = $1
    `
	err := r.db.Select(&users, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get users for team %s: %w", teamName, err)
	}
	return users, nil
}
