package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (h *Handler) GetBankingDetails(c *gin.Context) {
	profileId, err := getProfileId(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	details, err := h.bankingService.GetBankingDetails(profileId)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, details)
}

func (h *Handler) CreateRequisition(c *gin.Context) {
	profileId, err := getProfileId(c)
	req := &CreateRequisitionRequest{}
	err = c.BindJSON(req)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	accountDetails, err := h.bankingService.CreateRequisition(profileId, req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, accountDetails)
}

func (h *Handler) LinkAccountToProfile(c *gin.Context) {
	profileId, err := getProfileId(c)
	err = h.bankingService.LinkAccountToProfile(profileId)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) FilterTransactions(c *gin.Context) {
	profileId, err := getProfileId(c)
	params := &FilterTransactionsRequest{}
	err = c.BindQuery(params)
	if err != nil {
		return
	}
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	transactions, err := h.bankingService.FilterTransactions(profileId, params)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, transactions)
}

type CreateRequisitionResponse struct {
	Url string
}
type BankingDetailsResponse struct {
	BankAccounts []BankAccount
}
type BankAccount struct {
	OwnerName    string        `json:"ownerName"`
	Balance      float64       `json:"balance"`
	Iban         string        `json:"iban"`
	Transactions []Transaction `json:"transactions"`
}
type Transaction struct {
	Id                      uuid.UUID `json:"id"`
	TransactedAt            time.Time `json:"transactedAt"`
	Amount                  float64   `json:"amount"`
	CreditorName            string    `json:"creditorName"`
	CreditorIban            string    `json:"creditorIban"`
	DebtorName              string    `json:"debtorName"`
	DebtorIban              string    `json:"debtorIban"`
	BalanceAfterTransaction float64   `json:"balanceAfterTransaction"`
}
type CreateRequisitionRequest struct {
	InstitutionId  uuid.UUID `json:"institutionId"`
	MaxHistoryDays int64     `json:"maxHistoryDays"`
}
type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}
type FilterTransactionsRequest struct {
	TransactedBeforeTimestamp time.Time `form:"transactedBeforeTimestamp"`
	TransactedAfterTimestamp  time.Time `form:"transactedAfterTimestamp"`
	MinimumAmount             float64   `form:"minimumAmount"`
	MaximumAmount             float64   `form:"maximumAmount"`
	CreditorName              string    `form:"creditorName"`
	DebtorName                string    `form:"debtorName"`
	CreditorIban              string    `form:"creditorIban"`
	DebtorIban                string    `form:"debtorIban"`
}
