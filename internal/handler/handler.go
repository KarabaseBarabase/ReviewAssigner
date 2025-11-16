package handler

import (
	"ReviewAssigner/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	teamService *service.TeamService
	userService *service.UserService
	prService   *service.PRService
}

func NewHandler(
	teamService *service.TeamService,
	userService *service.UserService,
	prService *service.PRService,
) *Handler {
	return &Handler{
		teamService: teamService,
		userService: userService,
		prService:   prService,
	}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.GET("/health", h.healthCheck)

	router.POST("/team/add", h.addTeam)
	router.GET("/team/get", h.getTeam)
	router.POST("/team/:teamName/deactivate-users", h.deactivateUsers)

	router.POST("/users/setIsActive", h.setUserActive)
	router.GET("/users/getReview", h.getUserReviews)

	router.POST("/pullRequest/create", h.createPR)
	router.POST("/pullRequest/merge", h.mergePR)
	router.POST("/pullRequest/reassign", h.reassignReviewer)

	router.GET("/stats/user-assignments", h.getUserAssignmentsStats)
	router.GET("/stats/pr-metrics", h.getPRMetrics)
}
