package main

import (
	"os"
	"path/filepath"

	"github.com/QuoineFinancial/liquid-chain/cmd/node"
)

func main() {
	defaultRootDir := filepath.Join(os.Getenv("HOME"), ".liquid-chain")
	rootDir := defaultRootDir
	if rootDirEnv, ok := os.LookupEnv("ROOT_DIR"); ok {
		rootDir = rootDirEnv
	}

	// TODO: Get gasContractAddress from genesis file
	gasContractAddress := os.Getenv("GAS_CONTRACT_ADDRESS")
	liquidNode := node.New(rootDir, gasContractAddress)
	liquidNode.Execute()
}
