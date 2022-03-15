package node

import (
	"github.com/QuoineFinancial/liquid-chain/api"
	"github.com/QuoineFinancial/liquid-chain/consensus"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/libs/cli"
	tmNode "github.com/tendermint/tendermint/node"
)

// LiquidNode is the space where app and command lives
type LiquidNode struct {
	rootDir            string
	gasContractAddress string
	app                *consensus.App
	command            *cobra.Command
	tmNode             *tmNode.Node
	chainAPI           *api.API
}

// New returns new instance of Node
func New(rootDir string, gasContractAddress string) *LiquidNode {
	liquidNode := LiquidNode{
		rootDir:            rootDir,
		command:            commands.RootCmd,
		gasContractAddress: gasContractAddress,
	}
	liquidNode.addDefaultCommands()
	liquidNode.addStartNodeCommand()
	return &liquidNode
}

func (node *LiquidNode) addDefaultCommands() {
	node.command.AddCommand(
		commands.GenValidatorCmd,
		commands.InitFilesCmd,
		commands.ProbeUpnpCmd,
		commands.LiteCmd,
		commands.ReplayCmd,
		commands.ReplayConsoleCmd,
		commands.ResetAllCmd,
		commands.ResetPrivValidatorCmd,
		commands.ShowValidatorCmd,
		commands.TestnetFilesCmd,
		commands.ShowNodeIDCmd,
		commands.GenNodeKeyCmd,
		commands.VersionCmd,
	)

}

// Execute run the node.command base on user input
func (node *LiquidNode) Execute() {
	prefix := "TM"
	command := cli.PrepareBaseCmd(node.command, prefix, node.rootDir)
	if err := command.Execute(); err != nil {
		panic(err)
	}
}
