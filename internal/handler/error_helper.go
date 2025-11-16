package handler

import (
	"net/http"

	"ReviewAssigner/internal/errors"

	"github.com/gin-gonic/gin"
)

// ErrorResponse стандартный ответ с ошибкой
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// handleError обрабатывает ошибки и возвращает стандартизированный ответ
func handleError(c *gin.Context, err error) {
	if domainErr, ok := err.(*errors.Error); ok {
		status := getHTTPStatus(domainErr.Code)
		response := ErrorResponse{}
		response.Error.Code = domainErr.Code
		response.Error.Message = domainErr.Message
		c.JSON(status, response)
		return
	}

	// общая ошибка
	response := ErrorResponse{}
	response.Error.Code = "INTERNAL_ERROR"
	response.Error.Message = "Internal server error"
	c.JSON(http.StatusInternalServerError, response)
}

func getHTTPStatus(errorCode string) int {
	switch errorCode {
	case "NOT_FOUND":
		return http.StatusNotFound
	case "PR_EXISTS", "TEAM_EXISTS", "PR_MERGED", "NOT_ASSIGNED", "NO_CANDIDATE":
		return http.StatusConflict
	case "INVALID_REQUEST":
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func validateRequest(c *gin.Context, request interface{}) bool {
	if err := c.ShouldBindJSON(request); err != nil {
		response := ErrorResponse{}
		response.Error.Code = "INVALID_REQUEST"
		response.Error.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return false
	}
	return true
}

func validateRequiredParam(c *gin.Context, param, paramName string) bool {
	if param == "" {
		response := ErrorResponse{}
		response.Error.Code = "INVALID_REQUEST"
		response.Error.Message = paramName + " parameter is required"
		c.JSON(http.StatusBadRequest, response)
		return false
	}
	return true
}
