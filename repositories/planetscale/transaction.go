package planetscale

import (
	"budgetting/services/banking"
	"database/sql"
	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
)

func (r *Repository) FilterTransactions(profileId uuid.UUID) ([]*banking.Transaction, error) {
	var transactions []*banking.Transaction
	var filter string
	query := `SELECT * FROM transaction ` + filter
	res, err := r.client.Query(query, profileId)
	if err != nil {
		return nil, err
	}
	defer func(res *sql.Rows) {
		_ = res.Close()
	}(res)
	for res.Next() {
		t := &banking.Transaction{}
		var floatAmount float64
		var floatBalanceAfterTransaction float64
		scanErr := res.Scan(
			&t.Id,
			&t.BankAccountId,
			&t.ExternalId,
			&t.TransactedAt,
			&floatAmount,
			&t.CreditorName,
			&t.CreditorIban,
			&t.DebtorName,
			&t.DebtorIban,
			&floatBalanceAfterTransaction,
		)
		t.Amount.SetFloat64(floatAmount)
		t.BalanceAfterTransaction.SetFloat64(floatBalanceAfterTransaction)
		if scanErr != nil {
			return nil, scanErr
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
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
