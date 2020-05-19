package main

import (
	"context"
	"errors"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	// TODO: We can subscribe only for txs that contain request messages
	TxQuery = "tm.event = 'Tx'"
)

func runImpl(c *Context, l *Logger) error {
	l.Info(":rocket: Starting WebSocket subscriber")
	err := c.client.Start()
	if err != nil {
		return err
	}

	ctx, cxl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cxl()

	l.Info(":ear: Subscribing to events with query: %s...", TxQuery)
	eventChan, err := c.client.Subscribe(ctx, "", TxQuery)
	if err != nil {
		return err
	}

	for {
		select {
		case ev := <-eventChan:
			go handleTransaction(c, l, ev.Data.(tmtypes.EventDataTx).TxResult)
		}
	}
}

func runCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"r"},
		Short:   "Run the oracle process",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			keys, err := keybase.List()
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				return errors.New("No key available")
			}
			c.keys = make(chan keyring.Info, len(keys))
			for _, key := range keys {
				c.keys <- key
			}
			c.validator, err = sdk.ValAddressFromBech32(cfg.Validator)
			if err != nil {
				return err
			}
			err = sdk.VerifyAddressFormat(c.validator)
			if err != nil {
				return err
			}
			l := NewLogger()
			l.Info(":star: Creating HTTP client with node URI: %s", cfg.NodeURI)
			c.client, err = rpchttp.New(cfg.NodeURI, "/websocket")
			if err != nil {
				return err
			}
			return runImpl(c, l)
		},
	}
	return cmd
}
