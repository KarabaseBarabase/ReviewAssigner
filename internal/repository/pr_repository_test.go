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

func TestPRRepository_CreatePR(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewPRRepository(sqlxDB)

	pr := &models.PullRequest{
		PullRequestID:   "pr-1001",
		PullRequestName: "Add search",
		AuthorID:        "u1",
	}

	mock.ExpectExec(`INSERT INTO pull_requests`).
		WithArgs("pr-1001", "Add search", "u1", "OPEN").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreatePR(pr)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPRRepository_GetPRByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewPRRepository(sqlxDB)

	prRows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status", "created_at", "merged_at"}).
		AddRow("pr-1001", "Add search", "u1", "OPEN", time.Now(), nil)

	reviewerRows := sqlmock.NewRows([]string{"reviewer_id"}).
		AddRow("u2").
		AddRow("u3")

	mock.ExpectQuery(`SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at`).
		WithArgs("pr-1001").
		WillReturnRows(prRows)

	mock.ExpectQuery(`SELECT reviewer_id FROM pr_reviewers`).
		WithArgs("pr-1001").
		WillReturnRows(reviewerRows)

	pr, err := repo.GetPRByID("pr-1001")
	require.NoError(t, err)
	assert.Equal(t, "pr-1001", pr.PullRequestID)
	assert.Equal(t, "Add search", pr.PullRequestName)
	assert.Equal(t, "u1", pr.AuthorID)
	assert.Equal(t, "OPEN", pr.Status)
	assert.Len(t, pr.AssignedReviewers, 2)
	assert.Contains(t, pr.AssignedReviewers, "u2")
	assert.Contains(t, pr.AssignedReviewers, "u3")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPRRepository_MergePR(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewPRRepository(sqlxDB)

	mock.ExpectExec(`UPDATE pull_requests SET status = 'MERGED'`).
		WithArgs("pr-1001").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.MergePR("pr-1001")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPRRepository_AddPRReviewer(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewPRRepository(sqlxDB)

	mock.ExpectExec(`INSERT INTO pr_reviewers`).
		WithArgs("pr-1001", "u2").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddPRReviewer("pr-1001", "u2")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPRRepository_ReplacePRReviewer(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewPRRepository(sqlxDB)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE pr_reviewers SET is_active = false`).
		WithArgs("pr-1001", "u2").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`INSERT INTO pr_reviewers`).
		WithArgs("pr-1001", "u3").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = repo.ReplacePRReviewer("pr-1001", "u2", "u3")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPRRepository_GetAssignedPRs(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewPRRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status"}).
		AddRow("pr-1001", "Add search", "u1", "OPEN").
		AddRow("pr-1002", "Fix bug", "u3", "OPEN")

	mock.ExpectQuery(`SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status`).
		WithArgs("u2").
		WillReturnRows(rows)

	prs, err := repo.GetAssignedPRs("u2")
	require.NoError(t, err)
	assert.Len(t, prs, 2)
	assert.Equal(t, "pr-1001", prs[0].PullRequestID)
	assert.Equal(t, "pr-1002", prs[1].PullRequestID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPRRepository_GetUserAssignmentStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewPRRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{"reviewer_id", "assignment_count"}).
		AddRow("u1", 3).
		AddRow("u2", 5)

	mock.ExpectQuery(`SELECT reviewer_id, COUNT`).
		WillReturnRows(rows)

	stats, err := repo.GetUserAssignmentStats()
	require.NoError(t, err)
	assert.Equal(t, 3, stats["u1"])
	assert.Equal(t, 5, stats["u2"])
	assert.NoError(t, mock.ExpectationsWereMet())
}
