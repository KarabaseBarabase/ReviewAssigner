package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserAssignmentsStats godoc
// @Summary Статистика назначений по пользователям
// @Description Возвращает количество назначений PR на каждого пользователя
// @Tags statistics
// @Accept json
// @Produce json
// @Success 200 {object} StatsResponse "Статистика назначений"
// @Router /stats/user-assignments [get]
func (h *Handler) getUserAssignmentsStats(c *gin.Context) {
	stats, err := h.prService.GetUserAssignmentStats()
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, StatsResponse{UserAssignments: stats})
}

// GetPRMetrics godoc
// @Summary Метрики Pull Requests
// @Description Возвращает общую статистику по PR
// @Tags statistics
// @Accept json
// @Produce json
// @Success 200 {object} StatsResponse "Метрики PR"
// @Router /stats/pr-metrics [get]
func (h *Handler) getPRMetrics(c *gin.Context) {
	metrics, err := h.prService.GetPRMetrics()
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, StatsResponse{PRMetrics: metrics})
}
