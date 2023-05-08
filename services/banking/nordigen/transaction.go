package nordigen

import (
	"budgetting/services/banking"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (c *Client) GetTransactions(bankAccountId uuid.UUID) ([]*banking.Transaction, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/accounts/%s/transactions", c.NordigenUrl, bankAccountId.String()),
		nil,
	)
	transactionResponse := &TransactionResponse{}
	err = c.Execute(req, transactionResponse)
	if err != nil {
		return nil, err
	}
	var transactions []*banking.Transaction
	for _, transaction := range transactionResponse.Transactions.Booked {
		t := convertToTransaction(transaction)
		t.BankAccountId = bankAccountId
		transactions = append(transactions, t)
	}
	return transactions, nil
}

type TransactionResponse struct {
	Transactions Transactions `json:"transactions"`
}

type Transactions struct {
	Booked  []Transaction `json:"booked"`
	Pending []Transaction `json:"pending"`
}

type Transaction struct {
	TransactionId                          string                  `json:"transactionId"`
	BookingDate                            Time                    `json:"bookingDate"`
	TransactionAmount                      AmountCurrencyPair      `json:"transactionAmount"`
	DebtorName                             string                  `json:"debtorName"`
	DebtorAccount                          Account                 `json:"debtorAccount"`
	RemittanceInformationUnstructuredArray []string                `json:"remittanceInformationUnstructuredArray"`
	ProprietaryBankTransactionCode         string                  `json:"proprietaryBankTransactionCode"`
	BalanceAfterTransaction                BalanceAfterTransaction `json:"balanceAfterTransaction"`
	InternalTransactionId                  string                  `json:"internalTransactionId"`
	CreditorName                           string                  `json:"creditorName"`
	CreditorAccount                        Account                 `json:"creditorAccount"`
}

type BalanceAfterTransaction struct {
	BalanceAmount AmountCurrencyPair `json:"balanceAmount"`
	BalanceType   string             `json:"balanceType"`
}

func convertToTransaction(t Transaction) *banking.Transaction {
	return &banking.Transaction{
		Id:                      uuid.Nil,
		BankAccountId:           uuid.Nil,
		ExternalId:              t.TransactionId,
		TransactedAt:            time.Time(t.BookingDate),
		Amount:                  &t.TransactionAmount.Amount,
		CreditorName:            t.CreditorName,
		CreditorIban:            t.CreditorAccount.Iban,
		DebtorName:              t.DebtorName,
		DebtorIban:              t.DebtorAccount.Iban,
		BalanceAfterTransaction: &t.BalanceAfterTransaction.BalanceAmount.Amount,
	}
}
