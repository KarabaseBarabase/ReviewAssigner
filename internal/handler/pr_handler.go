// package handler

// import (
// 	"net/http"

// 	"ReviewAssigner/internal/models"

// 	"github.com/gin-gonic/gin"
// )

// type CreatePRRequest struct {
// 	PullRequestID   string `json:"pull_request_id" binding:"required" example:"pr-123"`
// 	PullRequestName string `json:"pull_request_name" binding:"required" example:"Fix login issue"`
// 	AuthorID        string `json:"author_id" binding:"required" example:"user-456"`
// }

// type MergePRRequest struct {
// 	PullRequestID string `json:"pull_request_id" binding:"required" example:"pr-123"`
// }

// type ReassignReviewerRequest struct {
// 	PullRequestID string `json:"pull_request_id" binding:"required" example:"pr-123"`
// 	OldUserID     string `json:"old_reviewer_id" binding:"required" example:"user-789"`
// }

// // CreatePR godoc
// // @Summary Создание Pull Request
// // @Description Создает новый PR и автоматически назначает ревьюеров
// // @Tags pull-requests
// // @Accept json
// // @Produce json
// // @Param request body CreatePRRequest true "Данные Pull Request" example:{"pull_request_id":"pr-123","pull_request_name":"Fix login issue","author_id":"user-456"}
// // @Success 201 {object} PRResponse "Созданный PR"
// // @Failure 400 {object} ErrorResponse "Ошибка валидации"
// // @Failure 409 {object} ErrorResponse "PR уже существует"
// // @Router /pullRequest/create [post]
// func (h *Handler) createPR(c *gin.Context) {
// 	var request CreatePRRequest

// 	if !validateRequest(c, &request) {
// 		return
// 	}

// 	pr := &models.PullRequest{
// 		PullRequestID:   request.PullRequestID,
// 		PullRequestName: request.PullRequestName,
// 		AuthorID:        request.AuthorID,
// 	}

// 	createdPR, err := h.prService.CreatePR(pr)
// 	if err != nil {
// 		handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusCreated, PRResponse{PR: createdPR})
// }

// // MergePR godoc
// // @Summary Merge Pull Request
// // @Description Помечает PR как мердженный
// // @Tags pull-requests
// // @Accept json
// // @Produce json
// // @Param request body MergePRRequest true "ID Pull Request" example:{"pull_request_id":"pr-123"}
// // @Success 200 {object} PRResponse "Обновленный PR"
// // @Failure 400 {object} ErrorResponse "Ошибка валидации"
// // @Failure 404 {object} ErrorResponse "PR не найден"
// // @Router /pullRequest/merge [post]
// func (h *Handler) mergePR(c *gin.Context) {
// 	var request MergePRRequest

// 	if !validateRequest(c, &request) {
// 		return
// 	}

// 	pr, err := h.prService.MergePR(request.PullRequestID)
// 	if err != nil {
// 		handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, PRResponse{PR: pr})
// }

// // ReassignReviewer godoc
// // @Summary Замена ревьюера
// // @Description Заменяет ревьюера на PR на другого активного участника команды
// // @Tags pull-requests
// // @Accept json
// // @Produce json
// // @Param request body ReassignReviewerRequest true "Данные для замены ревьюера" example:{"pull_request_id":"pr-123","old_reviewer_id":"user-789"}
// // @Success 200 {object} ReassignReviewerResponse "Результат замены"
// // @Failure 400 {object} ErrorResponse "Ошибка валидации"
// // @Failure 404 {object} ErrorResponse "PR или ревьюер не найден"
// // @Failure 409 {object} ErrorResponse "Ошибка замены ревьюера"
// // @Router /pullRequest/reassign [post]
// func (h *Handler) reassignReviewer(c *gin.Context) {
// 	var request ReassignReviewerRequest

// 	if !validateRequest(c, &request) {
// 		return
// 	}

// 	newReviewerID, err := h.prService.ReplaceReviewer(request.PullRequestID, request.OldUserID)
// 	if err != nil {
// 		handleError(c, err)
// 		return
// 	}

// 	pr, err := h.prService.GetPRByID(request.PullRequestID)
// 	if err != nil {
// 		handleError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, ReassignReviewerResponse{
// 		PR:         pr,
// 		ReplacedBy: newReviewerID,
// 	})
// }

package handler

import (
	"net/http"

	"ReviewAssigner/internal/models"

	"github.com/gin-gonic/gin"
)

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required" example:"pr-123"`
	PullRequestName string `json:"pull_request_name" binding:"required" example:"Fix login issue"`
	AuthorID        string `json:"author_id" binding:"required" example:"user-456"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required" example:"pr-123"`
}

type ReassignReviewerRequest struct {
	PullRequestID     string `json:"pull_request_id" binding:"required" example:"pr-123"`
	CurrentReviewerID string `json:"current_reviewer_id" binding:"required" example:"user-789"`
}

// CreatePR godoc
// @Summary Создание Pull Request
// @Description Создает новый PR и автоматически назначает ревьюеров
// @Tags pull-requests
// @Accept json
// @Produce json
// @Param request body CreatePRRequest true "Данные Pull Request" example:{"pull_request_id":"pr-123","pull_request_name":"Fix login issue","author_id":"user-456"}
// @Success 201 {object} PRResponse "Созданный PR"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 409 {object} ErrorResponse "PR уже существует"
// @Router /pullRequest/create [post]
func (h *Handler) createPR(c *gin.Context) {
	var request CreatePRRequest
	if !validateRequest(c, &request) {
		return
	}

	pr := &models.PullRequest{
		PullRequestID:   request.PullRequestID,
		PullRequestName: request.PullRequestName,
		AuthorID:        request.AuthorID,
	}

	createdPR, err := h.prService.CreatePR(pr)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, PRResponse{PR: createdPR})
}

// MergePR godoc
// @Summary Merge Pull Request
// @Description Помечает PR как мердженный
// @Tags pull-requests
// @Accept json
// @Produce json
// @Param request body MergePRRequest true "ID Pull Request" example:{"pull_request_id":"pr-123"}
// @Success 200 {object} PRResponse "Обновленный PR"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 404 {object} ErrorResponse "PR не найден"
// @Router /pullRequest/merge [post]
func (h *Handler) mergePR(c *gin.Context) {
	var request MergePRRequest
	if !validateRequest(c, &request) {
		return
	}

	pr, err := h.prService.MergePR(request.PullRequestID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, PRResponse{PR: pr})
}

// ReassignReviewer godoc
// @Summary Замена ревьюера
// @Description Заменяет ревьюера на PR на другого активного участника команды
// @Tags pull-requests
// @Accept json
// @Produce json
// @Param request body ReassignReviewerRequest true "Данные для замены ревьюера" example:{"pull_request_id":"pr-123","current_reviewer_id":"user-789"}
// @Success 200 {object} ReassignReviewerResponse "Результат замены"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 404 {object} ErrorResponse "PR или ревьюер не найден"
// @Failure 409 {object} ErrorResponse "Ошибка замены ревьюера"
// @Router /pullRequest/reassign [post]
func (h *Handler) reassignReviewer(c *gin.Context) {
	var request ReassignReviewerRequest
	if !validateRequest(c, &request) {
		return
	}

	newReviewerID, err := h.prService.ReplaceReviewer(request.PullRequestID, request.CurrentReviewerID)
	if err != nil {
		handleError(c, err)
		return
	}

	pr, err := h.prService.GetPRByID(request.PullRequestID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, ReassignReviewerResponse{
		PR:         pr,
		ReplacedBy: newReviewerID,
	})
}
