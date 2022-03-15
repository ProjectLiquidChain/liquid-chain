package node

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/QuoineFinancial/liquid-chain/api"
	"github.com/QuoineFinancial/liquid-chain/consensus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/config"
	tmFlags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmNode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

func (node *LiquidNode) newTendermintNode(config *config.Config, logger log.Logger) (*tmNode.Node, error) {
	node.app = consensus.NewApp(filepath.Join(config.DBDir(), "liquid"), node.gasContractAddress)
	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to load or gen node key %s: %w", config.NodeKeyFile(), err)
	}

	return tmNode.NewNode(config,
		privval.LoadOrGenFilePV(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(node.app),
		tmNode.DefaultGenesisDocProviderFunc(config),
		tmNode.DefaultDBProvider,
		tmNode.DefaultMetricsProvider(config.Instrumentation),
		logger.With("module", "node"),
	)
}

// ParseConfig parses the config file
func (node *LiquidNode) ParseConfig() (*config.Config, error) {
	conf := config.DefaultConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	conf.SetRoot(node.rootDir)
	config.EnsureRoot(node.rootDir)
	if err = conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}

	return conf, err
}

// StartTendermintNode starts the tmNode
func (node *LiquidNode) StartTendermintNode(conf *config.Config) error {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger, err := tmFlags.ParseLogLevel(conf.LogLevel, logger, config.DefaultLogLevel())
	if err != nil {
		return err
	}

	n, err := node.newTendermintNode(conf, logger)
	if err != nil {
		return fmt.Errorf("Failed to create node: %v", err)
	}
	node.tmNode = n

	// Stop upon receiving SIGTERM or CTRL-C.
	tmos.TrapSignal(logger, func() {
		node.Stop()
	})

	if err := n.Start(); err != nil {
		return fmt.Errorf("Failed to start node: %v", err)
	}
	logger.Info("Started node", "nodeInfo", n.Switch().NodeInfo())
	return nil
}

// Start runs the node with optional api given by flag --api
func (node *LiquidNode) Start(conf *config.Config, apiFlag bool) error {
	if err := node.StartTendermintNode(conf); err != nil {
		return err
	}

	if apiFlag {
		node.chainAPI = api.NewAPI(":5555", "tcp://localhost:26657", node.rootDir, *node.app.Meta, *node.app.State, *node.app.Chain)
		err := node.chainAPI.Serve()
		if err != nil {
			return err
		}
	}

	return nil
}

// Stop shutdowns the node
func (node *LiquidNode) Stop() {
	if node.chainAPI != nil {
		node.chainAPI.Close()
	}

	if node.tmNode.IsRunning() {
		_ = node.tmNode.Stop() // TODO: Properly handle error
	}
}

func (node *LiquidNode) addStartNodeCommand() {
	var apiFlag bool

	cmd := &cobra.Command{
		Use:   "start [--api]",
		Short: "Start the liquid node",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := node.ParseConfig()
			if err != nil {
				return fmt.Errorf("Failed to parse config: %v", err)
			}

			err = node.Start(conf, apiFlag)
			if err != nil {
				return err
			}

			// Run forever.
			select {}
		},
	}
	cmd.PersistentFlags().BoolVarP(&apiFlag, "api", "a", false, "start api")

	commands.AddNodeFlags(cmd)
	node.command.AddCommand(cmd)
}
