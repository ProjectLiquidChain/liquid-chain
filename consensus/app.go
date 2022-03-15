package consensus

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/constant"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/QuoineFinancial/liquid-chain/token"

	abciTypes "github.com/tendermint/tendermint/abci/types"
)

const (
	metaDBDir  = "meta.db"
	stateDBDir = "state.db"
	chainDBDir = "chain.db"
)

// App basic Tendermint base app
type App struct {
	abciTypes.BaseApplication

	Meta  *storage.MetaStorage
	State *storage.StateStorage
	Chain *storage.ChainStorage

	gasStation         gas.Station
	gasContractAddress string
}

// We use this code to communicate with Tendermint
// https://docs.tendermint.com/master/spec/abci/abci.html
const (
	ResponseCodeOK    = uint32(0)
	ResponseCodeNotOK = uint32(1)
)

func blockHashToAppHash(blockHash common.Hash) []byte {
	if blockHash == common.EmptyHash {
		return []byte{}
	}
	return blockHash.Bytes()
}

func appHashToBlockHash(appHash []byte) common.Hash {
	if len(appHash) == 0 {
		return common.EmptyHash
	}
	return common.BytesToHash(appHash)
}

// NewApp initializes a new app
func NewApp(dbDir string, gasContractAddress string) *App {
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.Mkdir(dbDir, os.ModePerm)
	}
	app := &App{
		Meta:               storage.NewMetaStorage(db.NewRocksDB(filepath.Join(dbDir, metaDBDir))),
		State:              storage.NewStateStorage(db.NewRocksDB(filepath.Join(dbDir, stateDBDir))),
		Chain:              storage.NewChainStorage(db.NewRocksDB(filepath.Join(dbDir, chainDBDir))),
		gasContractAddress: gasContractAddress,
	}
	app.SetGasStation(gas.NewFreeStation(app))
	return app
}

// BeginBlock begins new block
func (app *App) BeginBlock(req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	lastBlockHash := appHashToBlockHash(req.Header.AppHash)
	previousBlock := app.Chain.MustGetBlock(lastBlockHash)
	app.State.MustLoadState(previousBlock)
	app.Chain.ComposeBlock(previousBlock, req.Header.Time)
	for app.gasStation.Switch() {
	}
	return abciTypes.ResponseBeginBlock{}
}

// Info returns application chain info
func (app *App) Info(req abciTypes.RequestInfo) (resInfo abciTypes.ResponseInfo) {
	lastBlockHeight := app.Meta.LatestBlockHeight()
	lastBlockHash := app.Meta.BlockHeightToBlockHash(lastBlockHeight)
	return abciTypes.ResponseInfo{
		LastBlockHeight:  int64(lastBlockHeight),
		LastBlockAppHash: blockHashToAppHash(lastBlockHash),
	}
}

// CheckTx checks if submitted transaction is valid and can be passed to next step
func (app *App) CheckTx(req abciTypes.RequestCheckTx) abciTypes.ResponseCheckTx {
	if len(req.Tx) > constant.MaxTransactionSize {
		return abciTypes.ResponseCheckTx{
			Code: ResponseCodeNotOK,
			Log:  fmt.Sprintf("Transaction size exceed %dB", constant.MaxTransactionSize),
		}
	}

	tx, err := crypto.DecodeTransaction(req.GetTx())
	if err != nil {
		return abciTypes.ResponseCheckTx{
			Code: ResponseCodeNotOK,
			Log:  err.Error(),
		}
	}

	if err := app.validateTx(tx); err != nil {
		return abciTypes.ResponseCheckTx{
			Code: ResponseCodeNotOK,
			Log:  err.Error(),
		}
	}

	return abciTypes.ResponseCheckTx{Code: ResponseCodeOK}
}

//DeliverTx executes the submitted transaction
func (app *App) DeliverTx(req abciTypes.RequestDeliverTx) abciTypes.ResponseDeliverTx {
	tx, err := crypto.DecodeTransaction(req.GetTx())
	if err != nil {
		return abciTypes.ResponseDeliverTx{Code: ResponseCodeNotOK}
	}

	if err := app.validateTx(tx); err != nil {
		return abciTypes.ResponseDeliverTx{Code: ResponseCodeNotOK}
	}

	receipt, err := app.applyTransaction(tx)
	if err != nil {
		panic(err)
	}

	app.State.Commit()
	if err := app.Chain.AddTransactionWithReceipt(tx, receipt); err != nil {
		panic(err)
	}

	return abciTypes.ResponseDeliverTx{Code: ResponseCodeOK}
}

// Commit returns the state root of application storage. Called once all block processing is complete
func (app *App) Commit() abciTypes.ResponseCommit {
	blockHash := app.Chain.Commit(app.State.Commit())
	if err := app.Meta.StoreBlockMetas(app.Chain.CurrentBlock); err != nil {
		log.Println("unable to store index for block", blockHash)
	}
	return abciTypes.ResponseCommit{Data: blockHashToAppHash(blockHash)}
}

// SetGasStation active the gas station
func (app *App) SetGasStation(gasStation gas.Station) {
	app.gasStation = gasStation
}

// GetGasContractToken designated
func (app *App) GetGasContractToken() gas.Token {
	if len(app.gasContractAddress) > 0 {
		address, err := crypto.AddressFromString(app.gasContractAddress)
		if err != nil {
			panic(err)
		}
		contract, err := app.State.LoadAccount(address)
		if err != nil {
			panic(err)
		}
		if contract == nil {
			return nil
		}
		return token.NewToken(app.State, contract)
	}
	return nil
}
