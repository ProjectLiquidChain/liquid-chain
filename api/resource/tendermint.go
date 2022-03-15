package resource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	tmLog "github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/lite/proxy"
	"github.com/tendermint/tendermint/rpc/client"
	rpcHttp "github.com/tendermint/tendermint/rpc/client/http"
)

// TendermintAPI is client to interact with Tendermint RPC
type TendermintAPI = client.Client

const maxConnectionAttempt = 3
const apiInitDelay = 2 * time.Second

func readChainID(homeDir string) string {
	genesisPath := filepath.Join(homeDir, "/config/genesis.json")
	configFile, err := os.Open(genesisPath)
	if err != nil {
		panic("Unable to read genesis.json with error:\n" + err.Error())
	}
	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		panic("Invalid format of genesis.json")
	}
	var config struct {
		ChainID string `json:"chain_id"`
	}
	if err := json.Unmarshal(configBytes, &config); err != nil {
		panic("Could not read chain_id from genesis file")
	}
	return config.ChainID
}

// NewTendermintAPI returns new instance of TendermintAPI
func NewTendermintAPI(homeDir, nodeURL string) TendermintAPI {
	time.Sleep(apiInitDelay)
	chainID := readChainID(homeDir)
	logFileName := fmt.Sprintf("tendermint-api-%d.log", time.Now().Unix())
	logFilePath := filepath.Join(homeDir, logFileName)
	tendermintLoggerFile, _ := os.Create(logFilePath)
	defer tendermintLoggerFile.Close()
	logger := tmLog.NewTMLogger(tmLog.NewSyncWriter(tendermintLoggerFile))
	cacheSize := 10
	node, err := rpcHttp.New(nodeURL, "/websocket")
	if err != nil {
		tmos.Exit(err.Error())
	}

	cert, err := proxy.NewVerifier(chainID, homeDir, node, logger, cacheSize)
	cert.SetLogger(logger)
	return proxy.SecureClient(node, cert)
}
