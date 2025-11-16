// package handler

// import (
// 	"net/http"

// 	"ReviewAssigner/internal/models"

// 	"github.com/gin-gonic/gin"
// )

// type AddTeamRequest struct {
// 	TeamName string              `json:"team_name" binding:"required" example:"backend"`
// 	Members  []models.TeamMember `json:"members" binding:"required"`
// }

// // AddTeam godoc
// // @Summary Создание команды
// // @Description Создает новую команду с участниками
// // @Tags teams
// // @Accept json
// // @Produce json
// // @Param request body AddTeamRequest true "Данные команды" example:{"team_name":"backend","members":[{"user_id":"u1","username":"Alice","is_active":true},{"user_id":"u2","username":"Bob","is_active":true}]}
// // @Success 201 {object} TeamResponse "Созданная команда"
// // @Failure 400 {object} ErrorResponse "Ошибка валидации"
// // @Failure 409 {object} ErrorResponse "Команда уже существует"
// // @Router /team/add [post]
// func (h *Handler) addTeam(c *gin.Context) {
// 	var request AddTeamRequest

// 	if !validateRequest(c, &request) {
// 		return
// 	}

// 	team := &models.Team{
// 		TeamName: request.TeamName,
// 		Members:  request.Members,
// 	}

// 	if err := h.teamService.CreateTeam(team); err != nil {
// 		handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusCreated, TeamResponse{Team: team})
// }

// // GetTeam godoc
// // @Summary Получение информации о команде
// // @Description Возвращает информацию о команде и её участниках
// // @Tags teams
// // @Accept json
// // @Produce json
// // @Param team_name query string true "Название команды" example:backend
// // @Success 200 {object} models.Team "Информация о команде"
// // @Failure 400 {object} ErrorResponse "Ошибка валидации"
// // @Failure 404 {object} ErrorResponse "Команда не найдена"
// // @Router /team/get [get]
// func (h *Handler) getTeam(c *gin.Context) {
// 	teamName := c.Query("team_name")
// 	if !validateRequiredParam(c, teamName, "team_name") {
// 		return
// 	}

// 	team, err := h.teamService.GetTeam(teamName)
// 	if err != nil {
// 		handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, team)
// }

package handler

import (
	"net/http"

	"ReviewAssigner/internal/models"

	"github.com/gin-gonic/gin"
)

type AddTeamRequest struct {
	TeamName string              `json:"team_name" binding:"required" example:"backend"`
	Members  []models.TeamMember `json:"members" binding:"required"`
}

// AddTeam godoc
// @Summary Создание команды
// @Description Создает новую команду с участниками
// @Tags teams
// @Accept json
// @Produce json
// @Param request body AddTeamRequest true "Данные команды" example:{"team_name":"backend","members":[{"user_id":"u1","username":"Alice","is_active":true},{"user_id":"u2","username":"Bob","is_active":true}]}
// @Success 201 {object} TeamResponse "Созданная команда"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 409 {object} ErrorResponse "Команда уже существует"
// @Router /team/add [post]
func (h *Handler) addTeam(c *gin.Context) {
	var request AddTeamRequest

	if !validateRequest(c, &request) {
		return
	}

	team := &models.Team{
		TeamName: request.TeamName,
		Members:  request.Members,
	}

	if err := h.teamService.CreateTeam(team); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, TeamResponse{Team: team})
}

// GetTeam godoc
// @Summary Получение информации о команде
// @Description Возвращает информацию о команде и её участниках
// @Tags teams
// @Accept json
// @Produce json
// @Param team_name query string true "Название команды" example:backend
// @Success 200 {object} models.Team "Информация о команде"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 404 {object} ErrorResponse "Команда не найдена"
// @Router /team/get [get]
func (h *Handler) getTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if !validateRequiredParam(c, teamName, "team_name") {
		return
	}

	team, err := h.teamService.GetTeam(teamName)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, team)
}
