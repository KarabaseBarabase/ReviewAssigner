// package handler

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// type DeactivateUsersRequest struct {
// 	UserIDs []string `json:"user_ids" binding:"required" example:"user-123,user-456"`
// }

// // DeactivateUsers godoc
// // @Summary Массовая деактивация пользователей
// // @Description Деактивирует пользователей команды и переназначает открытые PR
// // @Tags teams
// // @Accept json
// // @Produce json
// // @Param teamName path string true "Название команды" example:backend
// // @Param request body DeactivateUsersRequest true "Список ID пользователей для деактивации" example:{"user_ids":["user-123","user-456"]}
// // @Success 200 {object} DeactivateUsersResponse "Результаты деактивации"
// // @Failure 400 {object} ErrorResponse "Ошибка валидации"
// // @Router /team/{teamName}/deactivate-users [post]
// func (h *Handler) deactivateUsers(c *gin.Context) {
// 	teamName := c.Param("teamName")
// 	if !validateRequiredParam(c, teamName, "team name") {
// 		return
// 	}

// 	var req DeactivateUsersRequest
// 	if !validateRequest(c, &req) {
// 		return
// 	}

// 	results := make(map[string]interface{})
// 	for _, userID := range req.UserIDs {
// 		user, err := h.userService.SetUserActive(userID, false)
// 		if err != nil {
// 			results[userID] = gin.H{
// 				"status":  "error",
// 				"message": err.Error(),
// 			}
// 		} else {
// 			results[userID] = gin.H{
// 				"status":  "success",
// 				"message": "deactivated successfully",
// 				"user":    user,
// 			}
// 		}
// 	}

// 	c.JSON(http.StatusOK, DeactivateUsersResponse{
// 		TeamName: teamName,
// 		Results:  results,
// 	})
// }

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeactivateUsersRequest struct {
	UserIDs []string `json:"user_ids" binding:"required" example:"user-123,user-456"`
}

// DeactivateUsers godoc
// @Summary Массовая деактивация пользователей
// @Description Деактивирует пользователей команды и переназначает открытые PR
// @Tags teams
// @Accept json
// @Produce json
// @Param teamName path string true "Название команды" example:backend
// @Param request body DeactivateUsersRequest true "Список ID пользователей для деактивации" example:{"user_ids":["user-123","user-456"]}
// @Success 200 {object} DeactivateUsersResponse "Результаты деактивации"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Router /team/{teamName}/deactivate-users [post]
func (h *Handler) deactivateUsers(c *gin.Context) {
	teamName := c.Param("teamName")
	if !validateRequiredParam(c, teamName, "team name") {
		return
	}

	var req DeactivateUsersRequest
	if !validateRequest(c, &req) {
		return
	}

	results := make(map[string]interface{})
	for _, userID := range req.UserIDs {
		user, err := h.userService.SetUserActive(userID, false)
		if err != nil {
			results[userID] = gin.H{
				"status":  "error",
				"message": err.Error(),
			}
		} else {
			results[userID] = gin.H{
				"status":  "success",
				"message": "deactivated successfully",
				"user":    user,
			}
		}
	}

	c.JSON(http.StatusOK, DeactivateUsersResponse{
		TeamName: teamName,
		Results:  results,
	})
}
