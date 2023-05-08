package banking

import (
	"budgetting/api/http/handler"
	"github.com/google/uuid"
	"time"
)

func (s *Service) CreateRequisition(profileId uuid.UUID, req *handler.CreateRequisitionRequest) (*handler.CreateRequisitionResponse, error) {
	var agreementId uuid.UUID
	if req.MaxHistoryDays != 90 {
		ai, err := s.NordigenClient.CreateEndUserAgreement(req.InstitutionId,
			req.MaxHistoryDays,
		)
		if err != nil {
			return nil, err
		}
		agreementId = ai
	}
	r, err := s.NordigenClient.CreateRequisition(req.InstitutionId, agreementId)
	if err != nil {
		return nil, err
	}
	err = s.repo.CreateRequisition(r)
	if err != nil {
		return nil, err
	}
	return &handler.CreateRequisitionResponse{Url: r.Link}, nil
}

type Requisition struct {
	Id            uuid.UUID
	ProfileId     uuid.UUID
	CreatedAt     time.Time
	AgreementId   uuid.UUID
	BankAccounts  []*BankAccount
	InstitutionId uuid.UUID
	Link          string
}
