package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Health check
// @Description Проверка работоспособности сервиса
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse "Статус сервиса"
// @Router /health [get]
func (h *Handler) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Service: "review-service",
	})
}
