package nordigen

import (
	"budgetting/services/banking"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (c *Client) CreateRequisition(institutionId uuid.UUID, agreementId uuid.UUID) (*banking.Requisition, error) {
	requisitionRequest := CreateRequisitionRequest{
		Redirect:         c.NordigenRedirectUrl,
		InstitutionId:    institutionId,
		AgreementId:      agreementId,
		UserLanguage:     "NL",
		AccountSelection: false,
	}
	requisitionRequest.Redirect = c.NordigenRedirectUrl
	body, err := json.Marshal(requisitionRequest)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/requisitions/", c.NordigenUrl),
		bytes.NewBuffer(body),
	)
	requisitionResponse := &RequisitionResponse{}
	err = c.Execute(req, requisitionResponse)
	if err != nil {
		return nil, err
	}

	return convertToRequisition(requisitionResponse), nil
}

func (c *Client) GetRequisition(requisitionId uuid.UUID) (*banking.Requisition, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/requisitions/%s", c.NordigenUrl, requisitionId),
		nil,
	)
	requisitionResponse := &RequisitionResponse{}
	err = c.Execute(req, requisitionResponse)
	if err != nil {
		return nil, err
	}

	return convertToRequisition(requisitionResponse), nil
}

type CreateRequisitionRequest struct {
	Redirect         string    `json:"redirect"`
	AgreementId      uuid.UUID `json:"agreement"`
	InstitutionId    uuid.UUID `json:"institution_id"`
	UserLanguage     string    `json:"user_language"`
	AccountSelection bool      `json:"account_selection"`
}

type RequisitionResponse struct {
	Id                uuid.UUID   `json:"id"`
	Created           time.Time   `json:"created"`
	Redirect          string      `json:"redirect"`
	Status            string      `json:"status"`
	InstitutionId     uuid.UUID   `json:"institution_id"`
	Agreement         uuid.UUID   `json:"agreement"`
	Reference         string      `json:"reference"`
	Accounts          []string    `json:"accounts"`
	UserLanguage      string      `json:"user_language"`
	Link              string      `json:"link"`
	Ssn               interface{} `json:"ssn"`
	AccountSelection  bool        `json:"account_selection"`
	RedirectImmediate bool        `json:"redirect_immediate"`
}

func convertToRequisition(r *RequisitionResponse) *banking.Requisition {
	return &banking.Requisition{
		Id:            r.Id,
		CreatedAt:     r.Created,
		AgreementId:   r.Agreement,
		InstitutionId: r.InstitutionId,
		Link:          r.Link,
	}
}
