package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}
