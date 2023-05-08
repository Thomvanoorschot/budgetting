package routes

import (
	"budgetting/api/http/handler"
	"github.com/gin-gonic/gin"
)

func (r *Router) SetupBankingRoutes(rg *gin.RouterGroup, h *handler.Handler) {
	g := rg.Group("/banking")
	g.Use(r.AuthMiddlewareFunc)

	g.GET("/details", h.GetBankingDetails)
	g.POST("/requisition", h.CreateRequisition)
	g.PUT("/link-account", h.LinkAccountToProfile)
	g.GET("/filter-transactions", h.FilterTransactions)
}
