package repository

import (
	"testing"
	"time"

	"ReviewAssigner/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateOrUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	user := &models.User{
		UserID:   "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}

	mock.ExpectExec(`INSERT INTO users`).
		WithArgs("u1", "Alice", "backend", true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateOrUpdateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_SetUserActive(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	mock.ExpectExec(`UPDATE users SET is_active`).
		WithArgs(false, "u1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.SetUserActive("u1", false)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_SetUserActive_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	mock.ExpectExec(`UPDATE users SET is_active`).
		WithArgs(false, "nonexistent").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.SetUserActive("nonexistent", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{"user_id", "username", "team_name", "is_active", "created_at", "updated_at"}).
		AddRow("u1", "Alice", "backend", true, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT user_id, username, team_name, is_active, created_at, updated_at`).
		WithArgs("u1").
		WillReturnRows(rows)

	user, err := repo.GetUserByID("u1")
	require.NoError(t, err)
	assert.Equal(t, "u1", user.UserID)
	assert.Equal(t, "Alice", user.Username)
	assert.Equal(t, "backend", user.TeamName)
	assert.True(t, user.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetActiveTeamMembers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{"user_id", "username", "team_name", "is_active", "created_at", "updated_at"}).
		AddRow("u2", "Bob", "backend", true, time.Now(), time.Now()).
		AddRow("u3", "Charlie", "backend", true, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT user_id, username, team_name, is_active, created_at, updated_at`).
		WithArgs("backend", "u1").
		WillReturnRows(rows)

	users, err := repo.GetActiveTeamMembers("backend", "u1")
	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "u2", users[0].UserID)
	assert.Equal(t, "u3", users[1].UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
