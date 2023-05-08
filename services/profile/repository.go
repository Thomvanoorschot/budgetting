package profile

import "github.com/google/uuid"

type Repository interface {
	CreateProfile(userId string) (*Profile, error)
	GetProfile(profileId uuid.UUID) (*Profile, error)
}
