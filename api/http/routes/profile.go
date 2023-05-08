package routes

import (
	"budgetting/api/http/handler"
	"github.com/gin-gonic/gin"
)

func (r *Router) SetupProfileRoutes(rg *gin.RouterGroup, h *handler.Handler) {
	g := rg.Group("/profile")
	g.Use(r.AuthMiddlewareFunc)

	g.GET("/", h.GetProfile)
	g.POST("/", h.CreateProfile)
}
