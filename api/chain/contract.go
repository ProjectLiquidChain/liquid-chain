package chain

import (
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// GetContractParams is params to GetAccount transaction
type GetContractParams struct {
	Address string `json:"address"`
}

// GetContractResult is result of GetAccount
type GetContractResult struct {
	Contract *abi.Contract `json:"contract"`
}

// GetContract gets contract from account state of given address
func (service *Service) GetContract(r *http.Request, params *GetContractParams, result *GetContractResult) error {
	service.syncLatestState()
	address, err := crypto.AddressFromString(params.Address)
	if err != nil {
		return err
	}
	account, err := service.state.GetAccount(address)
	if err != nil {
		return err
	}
	contract, err := account.GetContract()
	if err != nil {
		return err
	}
	result.Contract = contract
	return nil
}
