package repository

import (
	"fmt"

	"ReviewAssigner/internal/models"

	"github.com/jmoiron/sqlx"
)

// реализует UserRepositoryImpl интерфейс
type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateOrUpdateUser(user *models.User) error {
	query := `
		INSERT INTO users (user_id, username, team_name, is_active, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			username = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active,
			updated_at = NOW()
	`
	_, err := r.db.Exec(query, user.UserID, user.Username, user.TeamName, user.IsActive)
	return err
}

func (r *UserRepositoryImpl) SetUserActive(userID string, isActive bool) error {
	query := `UPDATE users SET is_active = $1, updated_at = NOW() WHERE user_id = $2`
	result, err := r.db.Exec(query, isActive, userID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepositoryImpl) GetUsersByTeam(teamName string) ([]models.User, error) {
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
	return users, err
}

func (r *UserRepositoryImpl) GetActiveTeamMembers(teamName string, excludeUserID string) ([]models.User, error) {
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
        AND is_active = true 
        AND user_id != $2
        ORDER BY RANDOM()
    `
	err := r.db.Select(&users, query, teamName, excludeUserID)
	return users, err
}

func (r *UserRepositoryImpl) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	query := `
        SELECT 
            user_id, 
            username, 
            team_name, 
            is_active, 
            created_at, 
            updated_at
        FROM users 
        WHERE user_id = $1
    `
	err := r.db.Get(&user, query, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}
