package planetscale

import (
	"budgetting/services/banking"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func (r *Repository) GetBankAccounts(profileId uuid.UUID) ([]*banking.BankAccount, error) {
	var bankAccounts []*banking.BankAccount
	query := `SELECT * FROM bankAccount BA where BA.profileId = ?`
	res, err := r.client.Query(query, profileId)
	if err != nil {
		return nil, err
	}
	defer func(res *sql.Rows) {
		_ = res.Close()
	}(res)
	for res.Next() {
		bankAccount := &banking.BankAccount{}
		var floatBalance float64
		scanErr := res.Scan(&bankAccount.Id,
			&bankAccount.Iban,
			&floatBalance,
			&bankAccount.ProfileId,
			&bankAccount.RequisitionId,
			&bankAccount.OwnerName)
		bankAccount.Balance.SetFloat64(floatBalance)
		if scanErr != nil {
			return nil, scanErr
		}
		bankAccounts = append(bankAccounts, bankAccount)
	}
	return bankAccounts, nil
}

func (r *Repository) GetBankAccountIds(profileId uuid.UUID) ([]uuid.UUID, error) {
	var bankAccountIds []uuid.UUID
	query := `SELECT BA.id FROM bankAccount BA where BA.profileId = ?`
	res, err := r.client.Query(query, profileId)
	if err != nil {
		return nil, err
	}
	defer func(res *sql.Rows) {
		_ = res.Close()
	}(res)
	for res.Next() {
		bankAccountId := uuid.Nil
		scanErr := res.Scan(&bankAccountId)
		if scanErr != nil {
			return nil, scanErr
		}
		bankAccountIds = append(bankAccountIds, bankAccountId)
	}
	return bankAccountIds, nil
}

func (r *Repository) CreateBankAccount(b *banking.BankAccount) error {
	query := `INSERT INTO bankAccount (id, profileId, requisitionId, balance, iban, ownerName) VALUES (?, ?, ?, ?, ?, ?)`
	floatBalance, _ := b.Balance.Float64()
	_, err := r.client.Exec(
		query,
		b.Id,
		b.ProfileId,
		b.RequisitionId,
		floatBalance,
		b.Iban,
		b.OwnerName,
	)
	if err != nil {
		return err
	}

	return nil
}
