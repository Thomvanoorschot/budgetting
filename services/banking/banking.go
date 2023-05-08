package banking

import (
	"budgetting/api/http/handler"
	"budgetting/utils"
	"github.com/google/uuid"
	"math/big"
)

type NordigenClient interface {
	GetTransactions(bankAccountId uuid.UUID) ([]*Transaction, error)
	GetAccountBalance(bankAccountId uuid.UUID) (*big.Rat, error)
	GetAccountOwner(bankAccountId uuid.UUID) (*AccountOwner, error)
	GetRequisition(requisitionId uuid.UUID) (*Requisition, error)
	CreateEndUserAgreement(institutionId uuid.UUID, maxHistoricalDays int64) (uuid.UUID, error)
	CreateRequisition(institutionId uuid.UUID, agreementId uuid.UUID) (*Requisition, error)
}

type Service struct {
	NordigenClient NordigenClient
	repo           Repository
}

func NewService(nordigenClient NordigenClient, repo Repository) *Service {
	return &Service{NordigenClient: nordigenClient, repo: repo}
}

func (s *Service) GetBankingDetails(profileId uuid.UUID) (*handler.BankingDetailsResponse, error) {
	bankAccountIds, err := s.repo.GetBankAccountIds(profileId)
	if err != nil {
		return nil, err
	}

	bankAccountCh := make(chan utils.Tuple2[*handler.BankAccount, error])
	respCh := make(chan *handler.BankingDetailsResponse)
	errCh := make(chan error)
	go func() {
		response := &handler.BankingDetailsResponse{}
		for {
			bankAccountResult := <-bankAccountCh
			if bankAccountResult.Err != nil {
				errCh <- bankAccountResult.Err
			}
			response.BankAccounts = append(response.BankAccounts, *bankAccountResult.Value)
			if len(response.BankAccounts) == len(bankAccountIds) {
				respCh <- response
				return
			}
		}
	}()
	for _, accountId := range bankAccountIds {
		utils.Async3(bankAccountCh, func() (*handler.BankAccount, error) {
			ba, bankAccountErr := s.getBankAccount(accountId)
			return ba, bankAccountErr
		})
	}
	select {
	case response := <-respCh:
		return response, nil
	case bankingAccountError := <-errCh:
		return nil, bankingAccountError
	}
}

func (s *Service) FilterTransactions(
	profileId uuid.UUID,
	params *handler.FilterTransactionsRequest,
) (*handler.TransactionsResponse, error) {
	//s.repo.FilterTransactions()

	return nil, nil
}

func (s *Service) getBankAccount(accountId uuid.UUID) (*handler.BankAccount, error) {
	transactionCh := utils.Async2(func() ([]*Transaction, error) {
		return s.NordigenClient.GetTransactions(accountId)
	})
	balanceCh := utils.Async2(func() (*big.Rat, error) {
		return s.NordigenClient.GetAccountBalance(accountId)
	})
	detailsCh := utils.Async2(func() (*AccountOwner, error) {
		return s.NordigenClient.GetAccountOwner(accountId)
	})
	transactionResult := <-transactionCh
	if transactionResult.Err != nil {
		return nil, transactionResult.Err
	}
	balanceResult := <-balanceCh
	if balanceResult.Err != nil {
		return nil, balanceResult.Err
	}
	detailsResult := <-detailsCh
	if detailsResult.Err != nil {
		return nil, detailsResult.Err
	}
	err := s.repo.CreateTransactions(transactionResult.Value)
	if err != nil {
		return nil, err
	}
	return convertToBankAccountResponse(
		transactionResult.Value,
		balanceResult.Value,
		detailsResult.Value,
	), nil
}

//func convertToBankAccountResponse(t *nordigen.Transactions,
//	b *nordigen.BalancesResponse,
//	a *nordigen.AccountDetails,
//) *handler.BankAccount {
//	var transactions []handler.Transaction
//	for _, transaction := range append(t.Booked, t.Pending...) {
//		transactions = append(transactions, handler.Transaction{
//			Id:                      transaction.TransactionId,
//			TransactedAt:            time.Time(transaction.BookingDate),
//			Amount:                  transaction.TransactionAmount.Amount.FloatString(2),
//			Currency:                transaction.TransactionAmount.Currency,
//			CreditorName:            transaction.CreditorName,
//			CreditorIban:            transaction.CreditorAccount.Iban,
//			DebtorName:              transaction.DebtorName,
//			DebtorIban:              transaction.DebtorAccount.Iban,
//			BalanceAfterTransaction: transaction.BalanceAfterTransaction.BalanceAmount.Amount.FloatString(2),
//		})
//	}
//
//	return &handler.BankAccount{
//		Transactions: transactions,
//		Balance:      b.Balances[0].BalanceAmount.Amount.FloatString(2),
//		Iban:         a.Iban,
//		OwnerName:    a.OwnerName,
//	}
//}
