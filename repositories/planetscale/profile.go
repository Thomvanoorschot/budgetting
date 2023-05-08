package planetscale

import (
	"budgetting/services/profile"
	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
)

func (r *Repository) GetProfile(profileId uuid.UUID) (*profile.Profile, error) {
	p := &profile.Profile{}
	query := `SELECT * FROM profile WHERE id = ?`
	err := r.client.QueryRow(query, profileId).Scan(&p.Id, &p.UserId)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repository) CreateProfile(userId string) (*profile.Profile, error) {
	p := &profile.Profile{
		Id:     uuid.New(),
		UserId: userId,
	}
	query := `INSERT INTO profile (id, userId) VALUES (?, ?)`
	_, err := r.client.Exec(query, p.Id, userId)
	if err != nil {
		return nil, err
	}
	return p, nil
}
