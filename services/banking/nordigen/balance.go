package nordigen

import (
	"fmt"
	"github.com/google/uuid"
	"math/big"
	"net/http"
)

func (c *Client) GetAccountBalance(bankAccountId uuid.UUID) (*big.Rat, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/accounts/%s/balances", c.NordigenUrl, bankAccountId.String()),
		nil,
	)
	balancesResponse := &BalancesResponse{}
	err = c.Execute(req, balancesResponse)
	if err != nil {
		return nil, err
	}

	return &balancesResponse.Balances[0].BalanceAmount.Amount, err
}

type BalancesResponse struct {
	Balances []Balance `json:"balances"`
}

type Balance struct {
	BalanceAmount AmountCurrencyPair `json:"balanceAmount"`
	BalanceType   string             `json:"balanceType"`
}
