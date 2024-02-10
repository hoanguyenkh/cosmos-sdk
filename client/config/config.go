package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
)

func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		ChainID:        "",
		KeyringBackend: "os",
		Output:         "text",
		Node:           "tcp://localhost:26657",
		BroadcastMode:  "sync",
	}
}

type ClientConfig struct {
	ChainID        string `mapstructure:"chain-id" json:"chain-id"`
	KeyringBackend string `mapstructure:"keyring-backend" json:"keyring-backend"`
	Output         string `mapstructure:"output" json:"output"`
	Node           string `mapstructure:"node" json:"node"`
	BroadcastMode  string `mapstructure:"broadcast-mode" json:"broadcast-mode"`
}

<<<<<<< HEAD
func (c *ClientConfig) SetChainID(chainID string) {
	c.ChainID = chainID
}

func (c *ClientConfig) SetKeyringBackend(keyringBackend string) {
	c.KeyringBackend = keyringBackend
}

func (c *ClientConfig) SetOutput(output string) {
	c.Output = output
}

func (c *ClientConfig) SetNode(node string) {
	c.Node = node
}

func (c *ClientConfig) SetBroadcastMode(broadcastMode string) {
	c.BroadcastMode = broadcastMode
}

// ReadFromClientConfig reads values from client.toml file and updates them in client Context
func ReadFromClientConfig(ctx client.Context) (client.Context, error) {
=======
// ReadFromClientConfig reads values from client.toml file and updates them in client.Context
// It uses CreateClientConfig internally with no custom template and custom config.
// Deprecated: use CreateClientConfig instead.
func ReadFromClientConfig(ctx client.Context) (client.Context, error) {
	return CreateClientConfig(ctx, "", nil)
}

// ReadDefaultValuesFromDefaultClientConfig reads default values from default client.toml file and updates them in client.Context
// The client.toml is then discarded.
func ReadDefaultValuesFromDefaultClientConfig(ctx client.Context, customClientTemplate string, customConfig interface{}) (client.Context, error) {
	prevHomeDir := ctx.HomeDir
	dir, err := os.MkdirTemp("", "simapp")
	if err != nil {
		return ctx, fmt.Errorf("couldn't create temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	ctx.HomeDir = dir
	ctx, err = CreateClientConfig(ctx, customClientTemplate, customConfig)
	if err != nil {
		return ctx, fmt.Errorf("couldn't create client config: %w", err)
	}

	ctx.HomeDir = prevHomeDir
	return ctx, nil
}

// CreateClientConfig reads the client.toml file and returns a new populated client.Context
// If the client.toml file does not exist, it creates one with default values.
// It takes a customClientTemplate and customConfig as input that can be used to overwrite the default config and enhance the client.toml file.
// The custom template/config must be both provided or be "" and nil.
func CreateClientConfig(ctx client.Context, customClientTemplate string, customConfig interface{}) (client.Context, error) {
>>>>>>> 72a56d993 (fix(simapp): fix default home (#19393))
	configPath := filepath.Join(ctx.HomeDir, "config")
	configFilePath := filepath.Join(configPath, "client.toml")
	conf := DefaultConfig()

	// when client.toml does not exist create and init with default values
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
			return ctx, fmt.Errorf("couldn't make client config: %v", err)
		}

		if ctx.ChainID != "" {
			conf.ChainID = ctx.ChainID // chain-id will be written to the client.toml while initiating the chain.
		}

		if err := writeConfigToFile(configFilePath, conf); err != nil {
			return ctx, fmt.Errorf("could not write client config to the file: %v", err)
		}
	}

	conf, err := getClientConfig(configPath, ctx.Viper)
	if err != nil {
		return ctx, fmt.Errorf("couldn't get client config: %v", err)
	}
	// we need to update KeyringDir field on Client Context first cause it is used in NewKeyringFromBackend
	ctx = ctx.WithOutputFormat(conf.Output).
		WithChainID(conf.ChainID).
		WithKeyringDir(ctx.HomeDir)

	keyring, err := client.NewKeyringFromBackend(ctx, conf.KeyringBackend)
	if err != nil {
		return ctx, fmt.Errorf("couldn't get keyring: %w", err)
	}

	ctx = ctx.WithKeyring(keyring)

	// https://github.com/cosmos/cosmos-sdk/issues/8986
	client, err := client.NewClientFromNode(conf.Node)
	if err != nil {
		return ctx, fmt.Errorf("couldn't get client from nodeURI: %v", err)
	}

	ctx = ctx.WithNodeURI(conf.Node).
		WithClient(client).
		WithBroadcastMode(conf.BroadcastMode)

	return ctx, nil
}
