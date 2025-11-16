package service

import (
	"fmt"
	"log/slog"
	"time"

	"ReviewAssigner/internal/errors"
	"ReviewAssigner/internal/models"
	"ReviewAssigner/internal/repository"
)

type PRService struct {
	prRepo        repository.PRRepository
	userRepo      repository.UserRepository
	reviewService *ReviewService
	logger        *slog.Logger
}

func NewPRService(
	prRepo repository.PRRepository,
	userRepo repository.UserRepository,
	reviewService *ReviewService,
	logger *slog.Logger,
) *PRService {
	if logger == nil {
		logger = slog.Default()
	}

	return &PRService{
		prRepo:        prRepo,
		userRepo:      userRepo,
		reviewService: reviewService,
		logger:        logger,
	}
}

func (s *PRService) GetPRByID(prID string) (*models.PullRequest, error) {
	s.logger.Debug("getting PR by ID", "pr_id", prID)

	pr, err := s.prRepo.GetPRByID(prID)
	if err != nil {
		s.logger.Error("failed to get PR by ID", "pr_id", prID, "error", err)
		return nil, errors.WrapError(errors.ErrPRNotFound, err)
	}

	s.logger.Debug("successfully retrieved PR", "pr_id", prID, "status", pr.Status)
	return pr, nil
}

func (s *PRService) CreatePR(pr *models.PullRequest) (*models.PullRequest, error) {
	start := time.Now()
	s.logger.Info("creating PR", "pr_id", pr.PullRequestID, "author_id", pr.AuthorID)

	// проверка существования
	exists, err := s.prRepo.PRExists(pr.PullRequestID)
	if err != nil {
		s.logger.Error("failed to check PR existence", "pr_id", pr.PullRequestID, "error", err)
		return nil, fmt.Errorf("failed to check PR existence: %w", err)
	}
	if exists {
		s.logger.Warn("PR already exists", "pr_id", pr.PullRequestID)
		return nil, errors.ErrPRExists
	}

	// проверка автора
	author, err := s.userRepo.GetUserByID(pr.AuthorID)
	if err != nil {
		s.logger.Error("author not found", "author_id", pr.AuthorID, "error", err)
		return nil, errors.ErrAuthorNotFound
	}

	pr.Status = "OPEN"
	now := time.Now()
	pr.CreatedAt = &now

	// создаём PR
	if err := s.prRepo.CreatePR(pr); err != nil {
		s.logger.Error("failed to create PR", "pr_id", pr.PullRequestID, "error", err)
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}

	// назначаем ревьюверов
	reviewers, err := s.reviewService.AssignReviewers(author.TeamName, pr.AuthorID, pr.PullRequestID)
	if err != nil {
		s.logger.Error("failed to assign reviewers, rolling back PR creation",
			"pr_id", pr.PullRequestID, "error", err)

		if delErr := s.prRepo.DeletePR(pr.PullRequestID); delErr != nil {
			s.logger.Error("failed to delete PR during rollback",
				"pr_id", pr.PullRequestID, "error", delErr)
		}
		return nil, fmt.Errorf("failed to assign reviewers: %w", err)
	}

	pr.AssignedReviewers = reviewers

	duration := time.Since(start)
	s.logger.Info("successfully created PR",
		"pr_id", pr.PullRequestID,
		"reviewers", reviewers,
		"duration", duration)

	if duration > 500*time.Millisecond {
		s.logger.Warn("CreatePR took too long",
			"pr_id", pr.PullRequestID,
			"duration", duration,
			"threshold", 500*time.Millisecond)
	}

	return pr, nil
}

func (s *PRService) MergePR(prID string) (*models.PullRequest, error) {
	s.logger.Info("merging PR", "pr_id", prID)

	pr, err := s.prRepo.GetPRByID(prID)
	if err != nil {
		s.logger.Error("PR not found for merge", "pr_id", prID, "error", err)
		return nil, errors.WrapError(errors.ErrPRNotFound, err)
	}

	if pr.Status == "MERGED" {
		s.logger.Info("PR already merged", "pr_id", prID)
		return pr, nil
	}

	if err := s.prRepo.MergePR(prID); err != nil {
		s.logger.Error("failed to merge PR", "pr_id", prID, "error", err)
		return nil, fmt.Errorf("failed to merge PR: %w", err)
	}

	s.logger.Info("successfully merged PR", "pr_id", prID)
	return s.prRepo.GetPRByID(prID)
}

func (s *PRService) ReplaceReviewer(prID, oldReviewerID string) (string, error) {
	s.logger.Info("replacing reviewer",
		"pr_id", prID,
		"old_reviewer_id", oldReviewerID)

	pr, err := s.prRepo.GetPRByID(prID)
	if err != nil {
		s.logger.Error("PR not found for reviewer replacement", "pr_id", prID, "error", err)
		return "", errors.WrapError(errors.ErrPRNotFound, err)
	}

	if pr.Status == "MERGED" {
		s.logger.Warn("attempted to replace reviewer on merged PR", "pr_id", prID)
		return "", errors.ErrPRMerged
	}

	assigned, err := s.prRepo.IsReviewerAssigned(prID, oldReviewerID)
	if err != nil {
		s.logger.Error("failed to check reviewer assignment",
			"pr_id", prID, "reviewer_id", oldReviewerID, "error", err)
		return "", fmt.Errorf("failed to check reviewer assignment: %w", err)
	}
	if !assigned {
		s.logger.Warn("reviewer not assigned to PR",
			"pr_id", prID, "reviewer_id", oldReviewerID)
		return "", errors.ErrNotAssigned
	}

	newReviewerID, err := s.reviewService.ReplaceReviewer(prID, oldReviewerID)
	if err != nil {
		s.logger.Error("failed to replace reviewer",
			"pr_id", prID, "old_reviewer_id", oldReviewerID, "error", err)
		return "", err
	}

	s.logger.Info("successfully replaced reviewer",
		"pr_id", prID,
		"old_reviewer_id", oldReviewerID,
		"new_reviewer_id", newReviewerID)

	return newReviewerID, nil
}

func (s *PRService) GetAssignedPRs(userID string) ([]models.PullRequestShort, error) {
	s.logger.Debug("getting assigned PRs for user", "user_id", userID)

	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		s.logger.Error("user not found", "user_id", userID, "error", err)
		return nil, errors.WrapError(errors.ErrUserNotFound, err)
	}

	prs, err := s.prRepo.GetAssignedPRs(userID)
	if err != nil {
		s.logger.Error("failed to get assigned PRs", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get assigned PRs: %w", err)
	}

	s.logger.Debug("retrieved assigned PRs", "user_id", userID, "count", len(prs))
	return prs, nil
}

func (s *PRService) GetUserAssignmentStats() (map[string]int, error) {
	s.logger.Debug("getting user assignment stats")

	stats, err := s.prRepo.GetUserAssignmentStats()
	if err != nil {
		s.logger.Error("failed to get user assignment stats", "error", err)
		return nil, fmt.Errorf("failed to get user assignment stats: %w", err)
	}

	s.logger.Debug("retrieved user assignment stats", "user_count", len(stats))
	return stats, nil
}

func (s *PRService) GetPRMetrics() (map[string]interface{}, error) {
	s.logger.Debug("getting PR metrics")

	metrics, err := s.prRepo.GetPRMetrics()
	if err != nil {
		s.logger.Error("failed to get PR metrics", "error", err)
		return nil, fmt.Errorf("failed to get PR metrics: %w", err)
	}

	s.logger.Debug("retrieved PR metrics", "metrics_count", len(metrics))
	return metrics, nil
}
