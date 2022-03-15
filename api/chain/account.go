package chain

import (
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

// GetAccountParams is params to GetAccount transaction
type GetAccountParams struct {
	Address string `json:"address"`
}

// GetAccountResult is result of GetAccount
type GetAccountResult struct {
	Account *storage.Account `json:"account"`
}

// GetAccount delivers transaction to blockchain
func (service *Service) GetAccount(r *http.Request, params *GetAccountParams, result *GetAccountResult) error {
	service.syncLatestState()
	address, err := crypto.AddressFromString(params.Address)
	if err != nil {
		return err
	}

	account, err := service.state.GetAccount(address)
	if err != nil {
		return err
	}
	result.Account = account
	return nil
}
