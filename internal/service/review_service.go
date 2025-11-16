package service

import (
	"fmt"
	"log/slog"
	"math/rand"

	"ReviewAssigner/internal/errors"
	"ReviewAssigner/internal/models"
	"ReviewAssigner/internal/repository"
)

type ReviewService struct {
	userRepo repository.UserRepository
	prRepo   repository.PRRepository
	logger   *slog.Logger
}

func NewReviewService(
	userRepo repository.UserRepository,
	prRepo repository.PRRepository,
	logger *slog.Logger,
) *ReviewService {
	if logger == nil {
		logger = slog.Default()
	}

	return &ReviewService{
		userRepo: userRepo,
		prRepo:   prRepo,
		logger:   logger,
	}
}

func (s *ReviewService) AssignReviewers(teamName, authorID, prID string) ([]string, error) {
	s.logger.Info("assigning reviewers",
		"team_name", teamName,
		"author_id", authorID,
		"pr_id", prID)

	if s.userRepo == nil {
		s.logger.Error("user repository is not initialized")
		return nil, fmt.Errorf("user repository is not initialized")
	}
	if s.prRepo == nil {
		s.logger.Error("PR repository is not initialized")
		return nil, fmt.Errorf("PR repository is not initialized")
	}

	candidates, err := s.userRepo.GetActiveTeamMembers(teamName, authorID)
	if err != nil {
		s.logger.Error("failed to get team members",
			"team_name", teamName, "author_id", authorID, "error", err)
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	s.logger.Debug("retrieved candidate reviewers",
		"team_name", teamName,
		"candidate_count", len(candidates))

	if len(candidates) == 0 {
		s.logger.Warn("no candidate reviewers available",
			"team_name", teamName, "author_id", authorID)
		return []string{}, nil
	}

	selected := s.selectRandomReviewers(candidates, 2)
	reviewerIDs := make([]string, 0, len(selected))

	for _, u := range selected {
		if err := s.prRepo.AddPRReviewer(prID, u.UserID); err != nil {
			s.logger.Error("failed to add PR reviewer",
				"pr_id", prID, "reviewer_id", u.UserID, "error", err)
			return nil, fmt.Errorf("failed to add reviewer %s: %w", u.UserID, err)
		}
		reviewerIDs = append(reviewerIDs, u.UserID)
	}

	s.logger.Info("successfully assigned reviewers",
		"pr_id", prID,
		"reviewers", reviewerIDs,
		"candidate_pool_size", len(candidates))

	return reviewerIDs, nil
}

func (s *ReviewService) ReplaceReviewer(prID, oldReviewerID string) (string, error) {
	s.logger.Info("replacing reviewer",
		"pr_id", prID,
		"old_reviewer_id", oldReviewerID)

	// информация о старом ревьювере
	oldReviewer, err := s.userRepo.GetUserByID(oldReviewerID)
	if err != nil {
		s.logger.Error("old reviewer not found",
			"reviewer_id", oldReviewerID, "error", err)
		return "", errors.WrapError(errors.ErrUserNotFound, err)
	}

	// текущие ревьюеры PR
	currentReviewers, err := s.prRepo.GetPRReviewers(prID)
	if err != nil {
		s.logger.Error("failed to get PR reviewers",
			"pr_id", prID, "error", err)
		return "", fmt.Errorf("failed to get PR reviewers: %w", err)
	}

	// кандидаты для замены
	candidates, err := s.userRepo.GetActiveTeamMembers(oldReviewer.TeamName, oldReviewerID)
	if err != nil {
		s.logger.Error("failed to get team members for replacement",
			"team_name", oldReviewer.TeamName, "error", err)
		return "", fmt.Errorf("failed to get team members: %w", err)
	}

	filteredCandidates := s.excludeUsers(candidates, currentReviewers)

	s.logger.Debug("reviewer replacement candidates",
		"pr_id", prID,
		"old_reviewer_id", oldReviewerID,
		"total_candidates", len(candidates),
		"filtered_candidates", len(filteredCandidates),
		"current_reviewers", currentReviewers)

	if len(filteredCandidates) == 0 {
		s.logger.Warn("no suitable candidates for reviewer replacement",
			"pr_id", prID, "old_reviewer_id", oldReviewerID)
		return "", errors.ErrNoCandidate
	}

	newReviewer := s.selectRandomReviewer(filteredCandidates)

	if err := s.prRepo.ReplacePRReviewer(prID, oldReviewerID, newReviewer.UserID); err != nil {
		s.logger.Error("failed to replace PR reviewer",
			"pr_id", prID,
			"old_reviewer_id", oldReviewerID,
			"new_reviewer_id", newReviewer.UserID,
			"error", err)
		return "", fmt.Errorf("failed to replace reviewer: %w", err)
	}

	s.logger.Info("successfully replaced reviewer",
		"pr_id", prID,
		"old_reviewer_id", oldReviewerID,
		"new_reviewer_id", newReviewer.UserID)

	return newReviewer.UserID, nil
}

// остальные методы остаются без изменений, только добавить логирование в selectRandomReviewers если нужно
func (s *ReviewService) selectRandomReviewers(candidates []models.User, max int) []models.User {
	if len(candidates) == 0 {
		return []models.User{}
	}
	if len(candidates) <= max {
		return candidates
	}
	shuffled := make([]models.User, len(candidates))
	copy(shuffled, candidates)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled[:max]
}

func (s *ReviewService) selectRandomReviewer(candidates []models.User) models.User {
	if len(candidates) == 0 {
		return models.User{}
	}
	return candidates[rand.Intn(len(candidates))]
}

func (s *ReviewService) excludeUsers(candidates []models.User, excludeIDs []string) []models.User {
	excludeMap := make(map[string]bool)
	for _, id := range excludeIDs {
		excludeMap[id] = true
	}
	var res []models.User
	for _, c := range candidates {
		if !excludeMap[c.UserID] {
			res = append(res, c)
		}
	}
	return res
}
