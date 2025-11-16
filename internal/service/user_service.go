package service

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"ReviewAssigner/internal/errors"
	"ReviewAssigner/internal/models"
	"ReviewAssigner/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
	prRepo   repository.PRRepository
	revSrv   *ReviewService
	logger   *slog.Logger
}

func NewUserService(
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
	prRepo repository.PRRepository,
	revSrv *ReviewService,
	logger *slog.Logger,
) *UserService {
	if logger == nil {
		logger = slog.Default()
	}

	return &UserService{
		userRepo: userRepo,
		teamRepo: teamRepo,
		prRepo:   prRepo,
		revSrv:   revSrv,
		logger:   logger,
	}
}

func (s *UserService) SetUserActive(userID string, isActive bool) (*models.User, error) {
	s.logger.Info("setting user active status",
		"user_id", userID, "is_active", isActive)

	if err := s.userRepo.SetUserActive(userID, isActive); err != nil {
		s.logger.Error("failed to set user active status",
			"user_id", userID, "is_active", isActive, "error", err)
		return nil, errors.WrapError(errors.ErrUserNotFound, err)
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		s.logger.Error("failed to get user after status change",
			"user_id", userID, "error", err)
		return nil, errors.WrapError(errors.ErrUserNotFound, err)
	}

	s.logger.Info("successfully changed user active status",
		"user_id", userID, "is_active", isActive)
	return user, nil
}

func (s *UserService) GetAssignedPRs(userID string) ([]models.PullRequestShort, error) {
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

func (s *UserService) BulkDeactivateUsers(teamName string, userIDs []string) (map[string]Reassignment, error) {
	start := time.Now()
	s.logger.Info("starting bulk deactivation",
		"team_name", teamName,
		"user_count", len(userIDs),
		"user_ids", userIDs)

	result := make(map[string]Reassignment)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// деактивируем пользователей
	for _, userID := range userIDs {
		if err := s.userRepo.SetUserActive(userID, false); err != nil {
			s.logger.Warn("failed to deactivate user",
				"user_id", userID, "error", err)
			mu.Lock()
			result[userID] = Reassignment{
				OldReviewer: userID,
				Success:     false,
				PRID:        "",
				NewReviewer: "",
			}
			mu.Unlock()
		} else {
			s.logger.Debug("successfully deactivated user", "user_id", userID)
			mu.Lock()
			result[userID] = Reassignment{
				OldReviewer: userID,
				Success:     true,
			}
			mu.Unlock()
		}
	}

	// для каждого пользователя получаем открытые PR и пытаемся заменить ревьюверов
	for _, userID := range userIDs {
		wg.Add(1)
		go func(uid string) {
			defer wg.Done()

			prs, err := s.prRepo.GetAssignedPRs(uid)
			if err != nil {
				s.logger.Warn("failed to get assigned PRs for user",
					"user_id", uid, "error", err)
				mu.Lock()
				result[uid] = Reassignment{
					OldReviewer: uid,
					Success:     false,
					PRID:        "",
					NewReviewer: "",
				}
				mu.Unlock()
				return
			}

			s.logger.Debug("found PRs assigned to user",
				"user_id", uid, "pr_count", len(prs))

			for _, pr := range prs {
				newReviewer, err := s.revSrv.ReplaceReviewer(pr.PullRequestID, uid)
				mu.Lock()
				if err != nil {
					s.logger.Error("failed to replace reviewer in PR",
						"pr_id", pr.PullRequestID,
						"old_reviewer_id", uid,
						"error", err)
					resultKey := uid + ":" + pr.PullRequestID
					result[resultKey] = Reassignment{
						OldReviewer: uid,
						PRID:        pr.PullRequestID,
						Success:     false,
					}
				} else {
					s.logger.Info("successfully replaced reviewer in PR",
						"pr_id", pr.PullRequestID,
						"old_reviewer_id", uid,
						"new_reviewer_id", newReviewer)
					resultKey := uid + ":" + pr.PullRequestID
					result[resultKey] = Reassignment{
						OldReviewer: uid,
						PRID:        pr.PullRequestID,
						NewReviewer: newReviewer,
						Success:     true,
					}
				}
				mu.Unlock()
			}
		}(userID)
	}

	wg.Wait()

	duration := time.Since(start)
	s.logger.Info("completed bulk deactivation",
		"team_name", teamName,
		"user_count", len(userIDs),
		"duration", duration,
		"result_count", len(result))

	if duration > 100*time.Millisecond {
		s.logger.Warn("BulkDeactivateUsers took too long",
			"duration", duration,
			"threshold", 100*time.Millisecond,
			"team_name", teamName,
			"user_count", len(userIDs))
	}

	return result, nil
}

type Reassignment struct {
	OldReviewer string `json:"old_reviewer"`
	NewReviewer string `json:"new_reviewer"`
	PRID        string `json:"pr_id"`
	Success     bool   `json:"success"`
}
