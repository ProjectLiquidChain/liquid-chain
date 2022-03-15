package chain

import (
	"github.com/QuoineFinancial/liquid-chain/api/resource"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

// Service is first service
type Service struct {
	tmAPI resource.TendermintAPI
	meta  *storage.MetaStorage
	state *storage.StateStorage
	block *storage.ChainStorage
}

// NewService returns new instance of Service
func NewService(
	tmAPI resource.TendermintAPI,
	meta *storage.MetaStorage,
	state *storage.StateStorage,
	block *storage.ChainStorage,
) *Service {
	return &Service{tmAPI, meta, state, block}
}

func (service *Service) syncStateAt(blockHeight uint64) {
	latestBlockHash := service.meta.BlockHeightToBlockHash(blockHeight)
	latestBlock := service.block.MustGetBlock(latestBlockHash)
	if err := service.state.LoadState(latestBlock); err != nil {
		panic(err)
	}
}

func (service *Service) syncLatestState() {
	service.syncStateAt(service.meta.LatestBlockHeight())
}
