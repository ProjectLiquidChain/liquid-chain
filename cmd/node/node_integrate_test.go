package node

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/QuoineFinancial/liquid-chain/api"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/util"
	"github.com/google/go-cmp/cmp"
	"github.com/tendermint/tendermint/config"
)

type testCase struct {
	name   string
	method string
	params string
	result string
}

type testServer struct {
	node *LiquidNode
}

const (
	blockchainTestName = "integration_test"
	gasContractAddress = "LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ7"
	SEED               = "0c61093a4983f5ba8cf83939efc6719e0c61093a4983f5ba8cf83939efc6719e"
)

func (ts *testServer) startNode() {
	conf := config.ResetTestRoot(blockchainTestName)
	fmt.Println("Init node config data...")

	ts.node = New(conf.RootDir, "")
	conf, err := ts.node.ParseConfig()
	if err != nil {
		panic(err)
	}
	conf.LogLevel = "error"
	conf.Consensus.CreateEmptyBlocks = false

	go func() {
		err := ts.node.StartTendermintNode(conf)
		if err != nil && err.Error() != "http: Server closed" {
			panic(err)
		}
	}()
	// Wait some time for server to ready
	time.Sleep(4 * time.Second)
}

// Please remember to call stopNode after done testing
func (ts *testServer) stopNode() {
	time.Sleep(2 * time.Second)

	ts.node.Stop()
	fmt.Println("Clean up node data")
	time.Sleep(500 * time.Millisecond)
	os.RemoveAll(ts.node.rootDir)

	time.Sleep(500 * time.Millisecond)
}

func loadPrivateKey(SEED string) ed25519.PrivateKey {
	hexSeed, err := hex.DecodeString(SEED)
	if err != nil {
		panic(err)
	}
	return ed25519.NewKeyFromSeed(hexSeed)
}

func getDeployLiquidTokenTx(t *testing.T, nonce uint64) string {
	seed := make([]byte, 32)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     nonce,
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	payload, err := util.BuildDeployTxPayload("../../test/testdata/liquid-token.wasm", "../../test/testdata/liquid-token-abi.json", "init", []string{"0"})
	if err != nil {
		t.Fatal(err)
	}
	deployTx := &crypto.Transaction{
		Version:   1,
		Sender:    &sender,
		Payload:   payload,
		Receiver:  crypto.EmptyAddress,
		GasLimit:  0,
		GasPrice:  1,
		Signature: nil,
	}
	dataToSign := crypto.GetSigHash(deployTx)
	deployTx.Signature = crypto.Sign(privateKey, dataToSign[:])
	rawTx, _ := deployTx.Encode()
	serializedTx := base64.StdEncoding.EncodeToString(rawTx)
	return serializedTx
}

func TestBroadcastTx(t *testing.T) {
	ts := &testServer{}
	defer ts.stopNode()
	ts.startNode()

	api := api.NewAPI(":5555", "tcp://localhost:26657", ts.node.rootDir, *ts.node.app.Meta, *ts.node.app.State, *ts.node.app.Chain)

	router := api.Router

	testcases := []testCase{
		{
			name:   "Broadcast Commit",
			method: "chain.BroadcastCommit",
			params: fmt.Sprintf(`{"rawTx": "%s"}`, getDeployLiquidTokenTx(t, 0)),
			result: `{"jsonrpc":"2.0","result":{"code":0,"log":"","hash":"420719772415f7902c8669678cfdf09b9a74c886652cae607786f6e6d2ee7cbc"},"id":1}`,
		},
		{
			name:   "Broadcast",
			method: "chain.Broadcast",
			params: fmt.Sprintf(`{"rawTx": "%s"}`, getDeployLiquidTokenTx(t, 1)),
			result: `{"jsonrpc":"2.0","result":{"code":0,"log":"","hash":"1f22b901c58ed95927c4c4d866289ecaf084d93d852044a795722ac2bd1a15a6"},"id":1}`,
		},
		{
			name:   "Broadcast Async",
			method: "chain.BroadcastAsync",
			params: fmt.Sprintf(`{"rawTx": "%s"}`, getDeployLiquidTokenTx(t, 2)),
			result: `{"jsonrpc":"2.0","result":{"code":0,"log":"","hash":"6de59e09494f6beb18928d3865a2293480aca86c4bea7846dd7d5e0b6cf6f4d1"},"id":1}`,
		},
	}

	for _, test := range testcases {
		response := httptest.NewRecorder()
		request, _ := makeRequest(test.method, test.params)
		router.ServeHTTP(response, request)
		result := readBody(response)
		fmt.Println(result)
		if diff := cmp.Diff(string(result), test.result); diff != "" {
			t.Errorf("%s: expect %s, got %s, diff: %s", test.name, test.result, result, diff)
		}
	}
}

func makeRequest(method string, params string) (*http.Request, error) {
	var body string
	if params == "" {
		body = fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "%s"}`, method)
	} else {
		body = fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "%s", "params": %s}`, method, params)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func readBody(res *httptest.ResponseRecorder) string {
	content, _ := ioutil.ReadAll(res.Body)
	stringResponse := strings.TrimSuffix(string(content), "\n")
	return string(stringResponse)
}
