package planetscale

import (
	"budgetting/config"
	"database/sql"
)

type Repository struct {
	client *sql.DB
}

func NewRepository(config *config.Config) (*Repository, error) {
	db, err := sql.Open("mysql", config.PlanetscaleDSN)
	if err != nil {
		return nil, err
	}
	return &Repository{client: db}, nil
}
