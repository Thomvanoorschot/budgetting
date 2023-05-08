package planetscale

import (
	"budgetting/services/banking"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func (r *Repository) GetRequisitionIds(profileId uuid.UUID) ([]uuid.UUID, error) {
	var requisitionIds []uuid.UUID
	query := `SELECT id FROM requisition where profileId = ?`
	res, err := r.client.Query(query, profileId)
	if err != nil {
		return nil, err
	}
	defer func(res *sql.Rows) {
		_ = res.Close()
	}(res)
	for res.Next() {
		requisitionId := uuid.Nil
		scanErr := res.Scan(&requisitionId)
		if scanErr != nil {
			return nil, scanErr
		}
		requisitionIds = append(requisitionIds, requisitionId)
	}
	return requisitionIds, nil
}
func (r *Repository) CreateRequisition(b *banking.Requisition) error {
	query := `INSERT INTO requisition (id, profileId, createdAt, agreementId, institutionId, link) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.client.Exec(query, b.Id, b.ProfileId, b.CreatedAt, b.AgreementId, b.InstitutionId, b.Link)
	if err != nil {
		return err
	}

	return nil
}
