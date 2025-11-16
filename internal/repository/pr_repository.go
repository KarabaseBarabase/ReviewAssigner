package repository

import (
	"fmt"

	"ReviewAssigner/internal/models"

	"github.com/jmoiron/sqlx"
)

// реализует PRRepository интерфейс
type PRRepositoryImpl struct {
	db *sqlx.DB
}

func NewPRRepository(db *sqlx.DB) *PRRepositoryImpl {
	return &PRRepositoryImpl{db: db}
}

func (r *PRRepositoryImpl) CreatePR(pr *models.PullRequest) error {
	query := `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := r.db.Exec(query, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, "OPEN")
	return err
}

func (r *PRRepositoryImpl) DeletePR(prID string) error {
	query := `DELETE FROM pull_requests WHERE pull_request_id = $1`
	_, err := r.db.Exec(query, prID)
	return err
}

func (r *PRRepositoryImpl) PRExists(prID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`
	err := r.db.Get(&exists, query, prID)
	return exists, err
}

func (r *PRRepositoryImpl) GetPRByID(prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	query := `
        SELECT 
            pull_request_id,
            pull_request_name, 
            author_id, 
            status,
            created_at,
            merged_at
        FROM pull_requests 
        WHERE pull_request_id = $1
    `
	err := r.db.Get(&pr, query, prID)
	if err != nil {
		return nil, fmt.Errorf("PR not found")
	}

	reviewers, err := r.GetPRReviewers(prID)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *PRRepositoryImpl) MergePR(prID string) error {
	query := `
		UPDATE pull_requests 
		SET status = 'MERGED', merged_at = NOW(), updated_at = NOW()
		WHERE pull_request_id = $1
	`
	_, err := r.db.Exec(query, prID)
	return err
}

func (r *PRRepositoryImpl) AddPRReviewer(prID, reviewerID string) error {
	// предполагаем, что на пару (pull_request_id, reviewer_id) есть уникальный индекс
	query := `
		INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at, is_active)
		VALUES ($1, $2, NOW(), true)
		ON CONFLICT (pull_request_id, reviewer_id)
		DO UPDATE SET is_active = true, replaced_at = NULL, assigned_at = NOW()
	`
	_, err := r.db.Exec(query, prID, reviewerID)
	return err
}

func (r *PRRepositoryImpl) ReplacePRReviewer(prID, oldReviewerID, newReviewerID string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	// если не коммит — откатим
	defer func() {
		_ = tx.Rollback()
	}()

	// деактивация старого ревьювера
	updateQuery := `
		UPDATE pr_reviewers 
		SET is_active = false, replaced_at = NOW()
		WHERE pull_request_id = $1 AND reviewer_id = $2 AND is_active = true
	`
	result, err := tx.Exec(updateQuery, prID, oldReviewerID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("reviewer not assigned to this PR")
	}

	// добавление или активация нового ревьювера
	insertQuery := `
		INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at, is_active)
		VALUES ($1, $2, NOW(), true)
		ON CONFLICT (pull_request_id, reviewer_id)
		DO UPDATE SET is_active = true, replaced_at = NULL, assigned_at = NOW()
	`
	_, err = tx.Exec(insertQuery, prID, newReviewerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PRRepositoryImpl) GetPRReviewers(prID string) ([]string, error) {
	var reviewers []string
	query := `
		SELECT reviewer_id FROM pr_reviewers 
		WHERE pull_request_id = $1 AND is_active = true
	`
	err := r.db.Select(&reviewers, query, prID)
	return reviewers, err
}

func (r *PRRepositoryImpl) GetAssignedPRs(userID string) ([]models.PullRequestShort, error) {
	var prs []models.PullRequestShort
	query := `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
		FROM pull_requests pr
		JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id
		WHERE prr.reviewer_id = $1 AND prr.is_active = true
		AND pr.status = 'OPEN'
	`
	err := r.db.Select(&prs, query, userID)
	return prs, err
}

func (r *PRRepositoryImpl) IsReviewerAssigned(prID, reviewerID string) (bool, error) {
	var assigned bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM pr_reviewers 
			WHERE pull_request_id = $1 AND reviewer_id = $2 AND is_active = true
		)
	`
	err := r.db.Get(&assigned, query, prID, reviewerID)
	return assigned, err
}

func (r *PRRepositoryImpl) GetUserAssignmentStats() (map[string]int, error) {
	type statsResult struct {
		ReviewerID      string `db:"reviewer_id"`
		AssignmentCount int    `db:"assignment_count"`
	}

	var results []statsResult
	query := `
		SELECT reviewer_id, COUNT(*) as assignment_count
		FROM pr_reviewers 
		WHERE is_active = true
		GROUP BY reviewer_id
	`

	err := r.db.Select(&results, query)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int)
	for _, result := range results {
		stats[result.ReviewerID] = result.AssignmentCount
	}

	return stats, nil
}

func (r *PRRepositoryImpl) GetPRMetrics() (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	var totalPRs int
	err := r.db.Get(&totalPRs, "SELECT COUNT(*) FROM pull_requests")
	if err != nil {
		return nil, err
	}
	metrics["total_prs"] = totalPRs

	var openPRs int
	err = r.db.Get(&openPRs, "SELECT COUNT(*) FROM pull_requests WHERE status = 'OPEN'")
	if err != nil {
		return nil, err
	}
	metrics["open_prs"] = openPRs

	var mergedPRs int
	err = r.db.Get(&mergedPRs, "SELECT COUNT(*) FROM pull_requests WHERE status = 'MERGED'")
	if err != nil {
		return nil, err
	}
	metrics["merged_prs"] = mergedPRs

	// ср.кол-во ревьюеров на PR (учитываем только активные записи)
	var avgReviewers float64
	err = r.db.Get(&avgReviewers, `
		SELECT COALESCE(AVG(reviewer_count), 0) 
		FROM (
			SELECT pull_request_id, COUNT(*) as reviewer_count 
			FROM pr_reviewers 
			WHERE is_active = true 
			GROUP BY pull_request_id
		) pr_counts
	`)
	if err != nil {
		return nil, err
	}
	metrics["avg_reviewers"] = avgReviewers

	return metrics, nil
}
