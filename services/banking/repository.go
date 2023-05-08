package banking

import "github.com/google/uuid"

type Repository interface {
	GetBankAccountIds(profileId uuid.UUID) ([]uuid.UUID, error)
	GetRequisitionIds(profileId uuid.UUID) ([]uuid.UUID, error)
	CreateRequisition(b *Requisition) error
	GetBankAccounts(profileId uuid.UUID) ([]*BankAccount, error)
	CreateBankAccount(b *BankAccount) error
	FilterTransactions(profileId uuid.UUID) ([]*Transaction, error)
	CreateTransactions(transactions []*Transaction) error
}
