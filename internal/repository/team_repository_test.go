package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeamRepository_CreateTeam(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewTeamRepository(sqlxDB)

	mock.ExpectExec(`INSERT INTO teams`).
		WithArgs("backend").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateTeam("backend")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_TeamExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewTeamRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs("backend").
		WillReturnRows(rows)

	exists, err := repo.TeamExists("backend")
	require.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetTeam(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewTeamRepository(sqlxDB)

	existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs("backend").
		WillReturnRows(existsRows)

	userRows := sqlmock.NewRows([]string{"user_id", "username", "team_name", "is_active", "created_at", "updated_at"}).
		AddRow("u1", "Alice", "backend", true, time.Now(), time.Now()).
		AddRow("u2", "Bob", "backend", true, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT user_id, username, team_name, is_active, created_at, updated_at`).
		WithArgs("backend").
		WillReturnRows(userRows)

	team, err := repo.GetTeam("backend")
	require.NoError(t, err)
	assert.Equal(t, "backend", team.TeamName)
	assert.Len(t, team.Members, 2)
	assert.Equal(t, "u1", team.Members[0].UserID)
	assert.Equal(t, "Alice", team.Members[0].Username)
	assert.Equal(t, "u2", team.Members[1].UserID)
	assert.Equal(t, "Bob", team.Members[1].Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepository_GetTeam_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewTeamRepository(sqlxDB)

	existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs("nonexistent").
		WillReturnRows(existsRows)

	team, err := repo.GetTeam("nonexistent")
	assert.Nil(t, team)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}
