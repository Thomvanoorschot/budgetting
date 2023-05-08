package handler

import (
	"context"
	"errors"
	jwtMiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BankingService interface {
	GetBankingDetails(profileId uuid.UUID) (*BankingDetailsResponse, error)
	CreateRequisition(
		profileId uuid.UUID,
		req *CreateRequisitionRequest,
	) (*CreateRequisitionResponse, error)
	LinkAccountToProfile(profileId uuid.UUID) error
	FilterTransactions(
		profileId uuid.UUID,
		params *FilterTransactionsRequest,
	) (*TransactionsResponse, error)
}

type ProfileService interface {
	GetProfile(profileId uuid.UUID) (*ProfileResponse, error)
	CreateProfile(userId string) error
}

type Handler struct {
	bankingService BankingService
	profileService ProfileService
}

type CustomClaims struct {
	ProfileId uuid.UUID `json:"profileId"`
}

func (c *CustomClaims) Validate(_ context.Context) error {
	return nil
}

func NewHandler(bankingService BankingService, profileService ProfileService) *Handler {
	return &Handler{bankingService: bankingService, profileService: profileService}
}

func getUserId(c *gin.Context) (string, error) {
	claims, ok := c.Request.Context().Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		return "", errors.New("could not parse claims from token")
	}
	return claims.RegisteredClaims.Subject, nil
}

func getProfileId(c *gin.Context) (uuid.UUID, error) {
	claims, ok := c.Request.Context().Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		return uuid.Nil, errors.New("could not parse claims from token")
	}
	customClaims, ok := claims.CustomClaims.(*CustomClaims)
	if !ok {
		return uuid.Nil, errors.New("could not parse custom claim")
	}
	return customClaims.ProfileId, nil
}
