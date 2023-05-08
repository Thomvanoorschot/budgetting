package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) GetProfile(c *gin.Context) {
	profileId, err := getProfileId(c)
	profile, err := h.profileService.GetProfile(profileId)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *Handler) CreateProfile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	err = h.profileService.CreateProfile(userId)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

type ProfileResponse struct {
	Id     uuid.UUID
	UserId string
}
