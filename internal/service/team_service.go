package service

import (
	"fmt"
	"log/slog"

	"ReviewAssigner/internal/errors"
	"ReviewAssigner/internal/models"
	"ReviewAssigner/internal/repository"
)

type TeamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewTeamService(
	teamRepo repository.TeamRepository,
	userRepo repository.UserRepository,
	logger *slog.Logger,
) *TeamService {
	if logger == nil {
		logger = slog.Default()
	}

	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *TeamService) CreateTeam(team *models.Team) error {
	s.logger.Info("creating team", "team_name", team.TeamName, "member_count", len(team.Members))

	exists, err := s.teamRepo.TeamExists(team.TeamName)
	if err != nil {
		s.logger.Error("failed to check team existence", "team_name", team.TeamName, "error", err)
		return fmt.Errorf("failed to check team existence: %w", err)
	}
	if exists {
		s.logger.Warn("team already exists", "team_name", team.TeamName)
		return errors.ErrTeamExists
	}

	if err := s.teamRepo.CreateTeam(team.TeamName); err != nil {
		s.logger.Error("failed to create team", "team_name", team.TeamName, "error", err)
		return fmt.Errorf("failed to create team: %w", err)
	}

	for _, member := range team.Members {
		user := &models.User{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: team.TeamName,
			IsActive: member.IsActive,
		}
		if err := s.userRepo.CreateOrUpdateUser(user); err != nil {
			s.logger.Error("failed to create/update team member",
				"team_name", team.TeamName,
				"user_id", member.UserID,
				"error", err)
			return fmt.Errorf("failed to create/update user %s: %w", member.UserID, err)
		}
	}

	s.logger.Info("successfully created team",
		"team_name", team.TeamName,
		"member_count", len(team.Members))
	return nil
}

func (s *TeamService) GetTeam(teamName string) (*models.Team, error) {
	s.logger.Debug("getting team", "team_name", teamName)

	team, err := s.teamRepo.GetTeam(teamName)
	if err != nil {
		s.logger.Error("team not found", "team_name", teamName, "error", err)
		return nil, errors.WrapError(errors.ErrTeamNotFound, err)
	}

	s.logger.Debug("successfully retrieved team",
		"team_name", teamName,
		"member_count", len(team.Members))
	return team, nil
}
