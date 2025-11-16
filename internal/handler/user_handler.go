// package handler

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// type SetUserActiveRequest struct {
// 	UserID   string `json:"user_id" binding:"required" example:"user-123"`
// 	IsActive bool   `json:"is_active" example:"false"`
// }

// // SetUserActive godoc
// // @Summary Установка активности пользователя
// // @Description Активирует или деактивирует пользователя
// // @Tags users
// // @Accept json
// // @Produce json
// // @Param request body SetUserActiveRequest true "Данные пользователя" example:{"user_id":"user-123","is_active":false}
// // @Success 200 {object} UserResponse "Обновленный пользователь"
// // @Failure 400 {object} ErrorResponse "Ошибка валидации"
// // @Failure 404 {object} ErrorResponse "Пользователь не найден"
// // @Router /users/setIsActive [post]
// func (h *Handler) setUserActive(c *gin.Context) {
// 	var request SetUserActiveRequest

// 	if !validateRequest(c, &request) {
// 		return
// 	}

// 	user, err := h.userService.SetUserActive(request.UserID, request.IsActive)
// 	if err != nil {
// 		handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, UserResponse{User: user})
// }

// // GetUserReviews godoc
// // @Summary Получение назначенных PR пользователя
// // @Description Возвращает список PR, назначенных на пользователя для ревью
// // @Tags users
// // @Accept json
// // @Produce json
// // @Param user_id query string true "ID пользователя" example:user-123
// // @Success 200 {object} UserPRsResponse "Список PR пользователя"
// // @Failure 400 {object} ErrorResponse "Ошибка валидации"
// // @Failure 404 {object} ErrorResponse "Пользователь не найден"
// // @Router /users/getReview [get]
// func (h *Handler) getUserReviews(c *gin.Context) {
// 	userID := c.Query("user_id")
// 	if !validateRequiredParam(c, userID, "user_id") {
// 		return
// 	}

// 	prs, err := h.userService.GetAssignedPRs(userID)
// 	if err != nil {
// 		handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, UserPRsResponse{
// 		UserID:       userID,
// 		PullRequests: prs,
// 	})
// }

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetUserActiveRequest struct {
	UserID   string `json:"user_id" binding:"required" example:"user-123"`
	IsActive bool   `json:"is_active" example:"false"`
}

// SetUserActive godoc
// @Summary Установка активности пользователя
// @Description Активирует или деактивирует пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param request body SetUserActiveRequest true "Данные пользователя" example:{"user_id":"user-123","is_active":false}
// @Success 200 {object} UserResponse "Обновленный пользователь"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Router /users/setIsActive [post]
func (h *Handler) setUserActive(c *gin.Context) {
	var request SetUserActiveRequest
	if !validateRequest(c, &request) {
		return
	}

	user, err := h.userService.SetUserActive(request.UserID, request.IsActive)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, UserResponse{User: user})
}

// GetUserReviews godoc
// @Summary Получение назначенных PR пользователя
// @Description Возвращает список PR, назначенных на пользователя для ревью
// @Tags users
// @Accept json
// @Produce json
// @Param user_id query string true "ID пользователя" example:user-123
// @Success 200 {object} UserPRsResponse "Список PR пользователя"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Router /users/getReview [get]
func (h *Handler) getUserReviews(c *gin.Context) {
	userID := c.Query("user_id")
	if !validateRequiredParam(c, userID, "user_id") {
		return
	}

	prs, err := h.userService.GetAssignedPRs(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, UserPRsResponse{
		UserID:       userID,
		PullRequests: prs,
	})
}
