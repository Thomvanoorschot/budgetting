package banking

import (
	"budgetting/api/http/handler"
	"budgetting/utils"
	"github.com/google/uuid"
	"math/big"
)

type BankAccount struct {
	Id            uuid.UUID      `json:"id"`
	ProfileId     uuid.UUID      `json:"profileId"`
	RequisitionId uuid.UUID      `json:"requisitionId"`
	Transactions  []*Transaction `json:"transactions"`
	Balance       *big.Rat       `json:"balance"`
	AccountOwner
}

type AccountOwner struct {
	Iban      string `json:"iban"`
	OwnerName string `json:"ownerName"`
}

func (s *Service) LinkAccountToProfile(profileId uuid.UUID) error {
	requisitionIds, err := s.repo.GetRequisitionIds(profileId)
	if err != nil {
		return err
	}
	bankAccounts, err := s.repo.GetBankAccounts(profileId)
	if err != nil {
		return err
	}
	bankAccountIds := map[uuid.UUID]bool{}
	processedRequisitionIds := map[uuid.UUID]bool{}
	for _, account := range bankAccounts {
		bankAccountIds[account.Id] = true
		processedRequisitionIds[account.RequisitionId] = true
	}
	for _, requisitionId := range requisitionIds {
		if processedRequisitionIds[requisitionId] {
			continue
		}
		requisition, err := s.NordigenClient.GetRequisition(requisitionId)
		if err != nil {
			return err
		}
		for _, requisitionBankAccount := range requisition.BankAccounts {
			bankAccountId := requisitionBankAccount.Id
			if bankAccountIds[bankAccountId] {
				continue
			}
			err := s.createBankAccount(requisitionId, profileId, bankAccountId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) createBankAccount(requisitionId uuid.UUID,
	profileId uuid.UUID,
	bankAccountId uuid.UUID) error {

	transactionCh := utils.Async2(func() ([]*Transaction, error) {
		return s.NordigenClient.GetTransactions(bankAccountId)
	})
	balanceCh := utils.Async2(func() (*big.Rat, error) {
		return s.NordigenClient.GetAccountBalance(bankAccountId)
	})
	detailsCh := utils.Async2(func() (*AccountOwner, error) {
		return s.NordigenClient.GetAccountOwner(bankAccountId)
	})
	transactionResult := <-transactionCh
	if transactionResult.Err != nil {
		return nil
	}
	balanceResult := <-balanceCh
	if balanceResult.Err != nil {
		return nil
	}
	detailsResult := <-detailsCh
	if detailsResult.Err != nil {
		return nil
	}

	bankAccount := &BankAccount{
		Id:            bankAccountId,
		Transactions:  transactionResult.Value,
		Balance:       balanceResult.Value,
		RequisitionId: requisitionId,
		ProfileId:     profileId,
		AccountOwner: AccountOwner{
			Iban:      detailsResult.Value.Iban,
			OwnerName: detailsResult.Value.OwnerName,
		},
	}
	err := s.repo.CreateBankAccount(bankAccount)
	if err != nil {
		return err
	}
	err = s.repo.CreateTransactions(transactionResult.Value)
	if err != nil {
		return err
	}
	return nil
}

func convertToBankAccountResponse(t []*Transaction,
	b *big.Rat,
	a *AccountOwner,
) *handler.BankAccount {
	var transactions []handler.Transaction
	for _, transaction := range t {
		transactions = append(transactions, convertToTransactionResponse(transaction))
	}
	floatBalance, _ := b.Float64()
	return &handler.BankAccount{
		Transactions: transactions,
		Balance:      floatBalance,
		Iban:         a.Iban,
		OwnerName:    a.OwnerName,
	}
}
