package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/healthz", h.healthz)
	r.GET("/notifications/:id", h.getNotification)
}

// healthz godoc
// @Summary Health check
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func (h Handler) healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "notification-service",
		"status":  "ok",
	})
}

// getNotification godoc
// @Summary Get notification detail
// @Tags notifications
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{}
// @Router /notifications/{id} [get]
func (h Handler) getNotification(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":      c.Param("id"),
		"message": "notification retrieval endpoint",
	})
}
