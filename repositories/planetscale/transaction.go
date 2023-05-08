package planetscale

import (
	"budgetting/services/banking"
	"database/sql"
	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
)

func (r *Repository) FilterTransactions(profileId uuid.UUID) ([]*banking.Transaction, error) {
	//p := &profile.Profile{}
	//query := `SELECT * FROM profile WHERE id = ?`

	//err := r.client.QueryRow(query, profileId).Scan(&p.Id, &p.UserId)
	//if err != nil {
	//	return nil, err
	//}
	//return p, nil
	return nil, nil
}

func (r *Repository) CreateTransactions(transactions []*banking.Transaction) error {
	tx, err := r.client.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	stmt, err := r.client.Prepare(
		`INSERT INTO transaction (
                         id,
                         bankAccountId,
                         externalId,
                         transactedAt,
                         amount,
                         creditorName,
                         creditorIban,
                         debtorName,
                         debtorIban,
                         balanceAfterTransaction
                         ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                         ON DUPLICATE KEY UPDATE 
                            bankAccountId=VALUES(bankAccountId),
                            externalId=VALUES(externalId),
                            transactedAt=VALUES(transactedAt),
                            amount=VALUES(amount),
                            creditorName=VALUES(creditorName),
                            creditorIban=VALUES(creditorIban),
                            debtorName=VALUES(debtorName),
                            debtorIban=VALUES(debtorIban),
                            balanceAfterTransaction=VALUES(balanceAfterTransaction)
			`,
	)
	if err != nil {
		return err
	}

	for _, transaction := range transactions {
		_, err = stmt.Exec(
			uuid.New(),
			transaction.BankAccountId,
			transaction.ExternalId,
			transaction.TransactedAt,
			transaction.Amount.FloatString(2),
			transaction.CreditorName,
			transaction.CreditorIban,
			transaction.DebtorName,
			transaction.DebtorIban,
			transaction.BalanceAfterTransaction.FloatString(2),
		)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
