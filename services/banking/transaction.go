package banking

import (
	"budgetting/api/http/handler"
	"github.com/google/uuid"
	"math/big"
	"time"
)

type Transaction struct {
	Id                      uuid.UUID `json:"id"`
	BankAccountId           uuid.UUID `json:"bankAccountId"`
	ExternalId              string    `json:"externalId"`
	TransactedAt            time.Time `json:"transactedAt"`
	Amount                  *big.Rat  `json:"amount"`
	CreditorName            string    `json:"creditorName"`
	CreditorIban            string    `json:"creditorIban"`
	DebtorName              string    `json:"debtorName"`
	DebtorIban              string    `json:"debtorIban"`
	BalanceAfterTransaction *big.Rat  `json:"balanceAfterTransaction"`
}

func convertToTransactionResponse(t *Transaction) handler.Transaction {
	floatAmount, _ := t.Amount.Float64()
	floatBalanceAfterTransaction, _ := t.BalanceAfterTransaction.Float64()
	return handler.Transaction{
		Id:                      t.Id,
		TransactedAt:            t.TransactedAt,
		Amount:                  floatAmount,
		CreditorName:            t.CreditorName,
		CreditorIban:            t.CreditorIban,
		DebtorName:              t.DebtorName,
		DebtorIban:              t.DebtorIban,
		BalanceAfterTransaction: floatBalanceAfterTransaction,
	}
}
