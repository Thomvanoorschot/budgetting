package nordigen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (c *Client) CreateEndUserAgreement(institutionId uuid.UUID, maxHistoricalDays int64) (uuid.UUID, error) {

	endUserAgreementReq := &CreateEndUserAgreementRequest{
		InstitutionId:      institutionId,
		MaxHistoricalDays:  fmt.Sprintf("%d", maxHistoricalDays),
		AccessValidForDays: "90",
		AccessScope:        []string{"balances", "details", "transactions"},
	}
	body, err := json.Marshal(endUserAgreementReq)
	if err != nil {
		return uuid.Nil, err
	}
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/agreements/enduser/", c.NordigenUrl),
		bytes.NewBuffer(body),
	)
	endUserAgreementResponse := &CreateEndUserAgreementResponse{}
	err = c.Execute(req, endUserAgreementResponse)
	if err != nil {
		return uuid.Nil, err
	}

	return endUserAgreementResponse.Id, nil
}

type CreateEndUserAgreementRequest struct {
	InstitutionId      uuid.UUID `json:"institution_id"`
	MaxHistoricalDays  string    `json:"max_historical_days"`
	AccessValidForDays string    `json:"access_valid_for_days"`
	AccessScope        []string  `json:"access_scope"`
}

type CreateEndUserAgreementResponse struct {
	Id                 uuid.UUID   `json:"id"`
	Created            time.Time   `json:"created"`
	InstitutionId      string      `json:"institution_id"`
	MaxHistoricalDays  int         `json:"max_historical_days"`
	AccessValidForDays int         `json:"access_valid_for_days"`
	AccessScope        []string    `json:"access_scope"`
	Accepted           interface{} `json:"accepted"`
}
