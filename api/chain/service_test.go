package chain

import (
	"crypto/ed25519"
	"encoding/base64"
	"os"
	"testing"

	"fmt"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/stretchr/testify/assert"
)

var testResourceInstance *testResource

func TestMain(m *testing.M) {
	testResourceInstance = newTestResource()
	testResourceInstance.seed()
	code := m.Run()
	testResourceInstance.tearDown()
	os.Exit(code)
}

func TestGetLatestBlock(t *testing.T) {
	var result BlockResult
	testResourceInstance.service.GetLatestBlock(nil, &LatestBlockParams{}, &result)

	fmt.Println(*result.Block)

	assert.Equal(t, block{
		Hash:            common.HexToHash("a3207d969ff4ba99d0194e1af35ef690ea5f71cfbf118bdccf62c77c35a95ea8"),
		Height:          4,
		Time:            4,
		Parent:          common.HexToHash("53637c9a17c42f0d5b181d9a78b23cd21d626fe61f2fd6fdd98df0771b1bea4b"),
		StateRoot:       common.HexToHash("9be9b60e9e8eec70c21d8c361263b501cc3edf0597c64ad8c2a79a2afaf0da41"),
		TransactionRoot: common.HexToHash("45b0cfc220ceec5b7c1c62c4d4193d38e4eba48e8815729ce75f9c0ab0e4c1c0"),
		ReceiptRoot:     common.HexToHash("45b0cfc220ceec5b7c1c62c4d4193d38e4eba48e8815729ce75f9c0ab0e4c1c0"),
		Transactions:    []transaction{},
		Receipts:        []receipt{},
	}, *result.Block)
}

func TestGetBlockByHeight(t *testing.T) {
	var result BlockResult
	testResourceInstance.service.GetBlockByHeight(nil, &BlockByHeightParams{
		Height: 2,
	}, &result)

	senders := make([]crypto.Address, 3)
	receivers := make([]crypto.Address, 3)

	for i := range senders {
		_, pk := testResourceInstance.getSenderWithNonce(byte(i), 0)
		senders[i] = crypto.AddressFromPubKey(pk.Public().(ed25519.PublicKey))
		receivers[i] = crypto.NewDeploymentAddress(senders[i], 0)
	}

	signatures := make([][]byte, 3)
	signatures[0], _ = base64.StdEncoding.DecodeString("OttqA4/C5Bk/04EMSXsBZ8U8bNWb4ErsBwStsdo4gDuV9kKEdb2Z/TEr9WQb100e7gj3g1meyKVinI2ZbjGcBg==")
	signatures[1], _ = base64.StdEncoding.DecodeString("qSZ39N1f1bOVE9Hhd9dog6iOHYbzXXNBr1fPjdUq5uOqrMbzh65gxPqpHNTZCHNgRRvC2QXLV+OReQqf4R9OCA==")
	signatures[2], _ = base64.StdEncoding.DecodeString("yyi3QL2b4IHPunS6m+cpc0dxD0k6Rx8frmQfHzc6hNt+gOZG4sBEesz92dGk2tnu/4dm7NNJKAQdBGYJoFPYDQ==")

	assert.Equal(t, block{
		Time:            2,
		Height:          2,
		Hash:            common.HexToHash("6ef55f78f482353d673e1429d30c967051ce8b7d47a3602e93d87e686bb5c7b0"),
		Parent:          common.HexToHash("acb376f46a530ef5c8d7702863d81e7c1f1b1a54f008cdbc7c380a8b8e13339b"),
		StateRoot:       common.HexToHash("3ec58cab3d13e0eaff2d4e06effb5445b9016324b38aef7f922157c987e5bbf7"),
		TransactionRoot: common.HexToHash("7c627e647b368cb1911bb95850210034927fb09394ff6be40548febe2db01b3d"),
		ReceiptRoot:     common.HexToHash("b54d0a78cdcdfba14bcaa39b7d7e2fccb56b594f4748fce3b6a55df9197e0758"),

		Transactions: []transaction{{
			Hash:        common.HexToHash("5e6552f82be4fe44e5f6915ca37ca2de24085da0cc83385040b68ace94b6d213"),
			Type:        "invoke",
			BlockHeight: 2,
			Version:     1,
			Sender:      senders[1],
			Nonce:       1,
			Receiver:    receivers[1],
			GasPrice:    1,
			GasLimit:    0,
			Signature:   signatures[1],
			Payload: call{
				Name: "mint",
				Args: []argument{{
					Type:  "uint64",
					Name:  "amount",
					Value: "1000",
				}},
			},
		}, {
			Hash:        common.HexToHash("b3fef26e5cb52f0681a06bb9c9ff78acb53ef79daa0c46fecbf1e212e9a67ddc"),
			Type:        "invoke",
			BlockHeight: 2,
			Version:     1,
			Sender:      senders[0],
			Nonce:       1,
			Receiver:    receivers[0],
			GasPrice:    1,
			GasLimit:    0,
			Signature:   signatures[0],
			Payload: call{
				Name: "mint",
				Args: []argument{{
					Type:  "uint64",
					Name:  "amount",
					Value: "1000",
				}},
			},
		}, {
			Hash:        common.HexToHash("c253c3e7f7ed6b393ba5a6e2ea73e4218abb3aa0eb5d36065a4367871341784b"),
			Type:        "invoke",
			BlockHeight: 2,
			Version:     1,
			Sender:      senders[2],
			Nonce:       1,
			Receiver:    receivers[2],
			GasPrice:    1,
			GasLimit:    0,
			Signature:   signatures[2],
			Payload: call{
				Name: "mint",
				Args: []argument{{
					Type:  "uint64",
					Name:  "amount",
					Value: "1000",
				}},
			},
		}},

		Receipts: []receipt{{
			Index:       0,
			Transaction: common.HexToHash("5e6552f82be4fe44e5f6915ca37ca2de24085da0cc83385040b68ace94b6d213"),
			Result:      "0",
			GasUsed:     0,
			Code:        0,
			Events: []call{{
				Contract: receivers[1].String(),
				Name:     "Mint",
				Args: []argument{{
					Type:  "address",
					Name:  "to",
					Value: senders[1].String(),
				}, {
					Type:  "uint64",
					Name:  "amount",
					Value: "1000",
				}},
			}},
			PostState: common.HexToHash("a6324a4017b35eb099ec4768d79e3a968f73510d1a5f38b42fce197d5fc021d7"),
		}, {
			Index:       1,
			Transaction: common.HexToHash("b3fef26e5cb52f0681a06bb9c9ff78acb53ef79daa0c46fecbf1e212e9a67ddc"),
			Result:      "0",
			GasUsed:     0,
			Code:        0,
			Events: []call{{
				Contract: receivers[0].String(),
				Name:     "Mint",
				Args: []argument{{
					Type:  "address",
					Name:  "to",
					Value: senders[0].String(),
				}, {
					Type:  "uint64",
					Name:  "amount",
					Value: "1000",
				}},
			}},
			PostState: common.HexToHash("3b214fa485b8125b9221fac1ae38dc24b259759772af8800044dc67caef5979c"),
		}, {
			Index:       2,
			Transaction: common.HexToHash("c253c3e7f7ed6b393ba5a6e2ea73e4218abb3aa0eb5d36065a4367871341784b"),
			Result:      "0",
			GasUsed:     0,
			Code:        0,
			Events: []call{{
				Contract: receivers[2].String(),
				Name:     "Mint",
				Args: []argument{{
					Type:  "address",
					Name:  "to",
					Value: senders[2].String(),
				}, {
					Type:  "uint64",
					Name:  "amount",
					Value: "1000",
				}},
			}},
			PostState: common.HexToHash("3ec58cab3d13e0eaff2d4e06effb5445b9016324b38aef7f922157c987e5bbf7"),
		}},
	}, *result.Block)
}

func TestGetTransaction(t *testing.T) {
	var result GetTransactionResult
	testResourceInstance.service.GetTransaction(nil, &GetTransactionParams{
		Hash: "b3fef26e5cb52f0681a06bb9c9ff78acb53ef79daa0c46fecbf1e212e9a67ddc",
	}, &result)

	sender, _ := crypto.AddressFromString("LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT")
	receiver, _ := crypto.AddressFromString("LBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH")
	signature, _ := base64.StdEncoding.DecodeString("OttqA4/C5Bk/04EMSXsBZ8U8bNWb4ErsBwStsdo4gDuV9kKEdb2Z/TEr9WQb100e7gj3g1meyKVinI2ZbjGcBg==")

	assert.Equal(t, transaction{
		Hash:        common.HexToHash("b3fef26e5cb52f0681a06bb9c9ff78acb53ef79daa0c46fecbf1e212e9a67ddc"),
		Type:        "invoke",
		BlockHeight: 2,
		Version:     1,
		Sender:      sender,
		Nonce:       1,
		Receiver:    receiver,
		GasPrice:    1,
		GasLimit:    0,
		Signature:   signature,
		Payload: call{
			Name: "mint",
			Args: []argument{{
				Type:  "uint64",
				Name:  "amount",
				Value: "1000",
			}},
		},
	}, *result.Transaction)

	assert.Equal(t, receipt{
		Index:       1,
		Transaction: common.HexToHash("b3fef26e5cb52f0681a06bb9c9ff78acb53ef79daa0c46fecbf1e212e9a67ddc"),
		Result:      "0",
		GasUsed:     0,
		Code:        0,
		Events: []call{{
			Name:     "Mint",
			Contract: receiver.String(),
			Args: []argument{{
				Type:  "address",
				Name:  "to",
				Value: "LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT",
			}, {
				Type:  "uint64",
				Name:  "amount",
				Value: "1000",
			}},
		}},
		PostState: common.HexToHash("3b214fa485b8125b9221fac1ae38dc24b259759772af8800044dc67caef5979c"),
	}, *result.Receipt)
}

func newUint64(value uint64) *uint64 {
	return &value
}

func TestCall(t *testing.T) {
	tests := []struct {
		name    string
		params  CallParams
		result  CallResult
		wantErr bool
	}{{
		name: "valid",
		params: CallParams{
			Height:  nil,
			Address: "LBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH",
			Method:  "get_balance",
			Args:    []string{"LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT"},
		},
		result: CallResult{
			Result: "3e8",
			Code:   crypto.ReceiptCodeOK,
			Events: []*call{},
		},
		wantErr: false,
	}, {
		name: "invalid address",
		params: CallParams{
			Height:  newUint64(1),
			Address: "invalid_address",
		},
		wantErr: true,
	}, {
		name: "call nil contract",
		params: CallParams{
			Height:  newUint64(1),
			Address: "LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53",
		},
		wantErr: true,
	}, {
		name: "call not a contract",
		params: CallParams{
			Height:  newUint64(1),
			Address: "LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT",
		},
		wantErr: true,
	}, {
		name: "invalid function",
		params: CallParams{
			Height:  newUint64(1),
			Address: "LBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH",
			Method:  "invalid_function",
		},
		wantErr: true,
	}, {
		name: "invalid params",
		params: CallParams{
			Height:  newUint64(1),
			Address: "LBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH",
			Method:  "get_balance",
			Args:    []string{},
		},
		wantErr: true,
	}, {
		name: "ignite with events",
		params: CallParams{
			Address: "LA3K6XGDQXAZN6J22J5VCEFIU25PE4BEZRZE5K76WDGUIRV3HLKJALPV",
			Method:  "say",
			Args:    []string{"1"},
		},
		result: CallResult{
			Result: "1",
			Code:   0,
			Events: []*call{{
				Contract: "LA3K6XGDQXAZN6J22J5VCEFIU25PE4BEZRZE5K76WDGUIRV3HLKJALPV",
				Name:     "Say",
				Args: []argument{{
					Type:  "uint8[]",
					Name:  "message",
					Value: "Q2hlY2tpbmc=",
				}},
			}},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result CallResult
			if err := testResourceInstance.service.Call(nil, &tt.params, &result); (err != nil) != tt.wantErr {
				t.Errorf("Service.Call() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestGetAccount(t *testing.T) {
	sender, _ := crypto.AddressFromString("LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT")

	tests := []struct {
		name    string
		params  GetAccountParams
		result  GetAccountResult
		wantErr bool
	}{{
		name: "valid",
		params: GetAccountParams{
			Address: "LBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH",
		},
		result: GetAccountResult{
			Account: &storage.Account{
				Nonce:        0,
				Creator:      sender,
				StorageHash:  common.Hash{0x29, 0xc3, 0x5c, 0xda, 0xdc, 0x63, 0x49, 0xf, 0xb9, 0x2d, 0xdf, 0x18, 0x80, 0xc0, 0xb2, 0x98, 0x29, 0xb2, 0xab, 0x82, 0x1d, 0xf9, 0x18, 0x58, 0x2f, 0xef, 0x98, 0x9, 0x5, 0xf1, 0x88, 0x5c},
				ContractHash: common.Hash{0xd8, 0x9a, 0xb7, 0x4c, 0xc7, 0xf9, 0x5c, 0x3, 0xd5, 0x7d, 0xc6, 0x76, 0xee, 0xeb, 0x9d, 0xfc, 0x78, 0x15, 0xde, 0xe8, 0xc0, 0x5d, 0x7b, 0x2a, 0xe2, 0x8b, 0x7, 0xee, 0x5f, 0x6a, 0xa1, 0x4},
			},
		},
		wantErr: false,
	}, {
		name: "invalid address",
		params: GetAccountParams{
			Address: "MBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH",
		},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result GetAccountResult
			err := testResourceInstance.service.GetAccount(nil, &tt.params, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				assert.Equal(t, tt.result.Account.Nonce, result.Account.Nonce, "Nonce")
				assert.Equal(t, tt.result.Account.Creator, result.Account.Creator, "Creator")
				assert.Equal(t, tt.result.Account.StorageHash, result.Account.StorageHash, "StorageHash")
				assert.Equal(t, tt.result.Account.ContractHash, result.Account.ContractHash, "ContractHash")
			}
		})
	}
}
