package routes

import (
	"budgetting/api/http/handler"
	"github.com/gin-gonic/gin"
)

func (r *Router) SetupHealthRoutes(rg *gin.RouterGroup, handler *handler.Handler) {
	rg.GET("/health", handler.GetHealth)
}
