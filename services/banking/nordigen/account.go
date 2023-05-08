package nordigen

import (
	"budgetting/services/banking"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func (c *Client) GetAccountOwner(bankAccountId uuid.UUID) (*banking.AccountOwner, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/accounts/%s/details", c.NordigenUrl, bankAccountId.String()),
		nil,
	)
	accountDetailResponse := &AccountDetailsResponse{}
	err = c.Execute(req, accountDetailResponse)

	if err != nil {
		return nil, err
	}

	return &banking.AccountOwner{
		Iban:      accountDetailResponse.AccountDetails.Iban,
		OwnerName: accountDetailResponse.AccountDetails.OwnerName,
	}, err
}

type AccountDetailsResponse struct {
	AccountDetails AccountDetails `json:"account"`
}
type AccountDetails struct {
	Account
	Currency  string `json:"currency"`
	OwnerName string `json:"ownerName"`
}
type Account struct {
	Iban string `json:"iban"`
}
