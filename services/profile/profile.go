package profile

import (
	"budgetting/api/http/handler"
	"budgetting/services/profile/auth0"
	"github.com/google/uuid"
)

type Service struct {
	repo        Repository
	auth0Client *auth0.Client
}

func NewService(repo Repository, auth0Client *auth0.Client) *Service {
	return &Service{repo: repo, auth0Client: auth0Client}
}

func (s *Service) GetProfile(profileId uuid.UUID) (*handler.ProfileResponse, error) {
	profile, err := s.repo.GetProfile(profileId)
	if err != nil {
		return nil, err
	}
	return &handler.ProfileResponse{Id: profile.Id, UserId: profile.UserId}, nil
}

func (s *Service) CreateProfile(userId string) error {
	p, err := s.repo.CreateProfile(userId)
	if err != nil {
		return err
	}
	return s.auth0Client.UpdateAppMetadata(userId, p.Id)
}

type Profile struct {
	Id            uuid.UUID
	UserId        string
	RequisitionId uuid.UUID
}
