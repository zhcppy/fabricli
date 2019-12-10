package event

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhcppy/fabricli/api"
	"github.com/zhcppy/fabricli/cmd"
)

func NewCmd() *cobra.Command {
	var eventCmd = &cobra.Command{
		Use:   "listen",
		Short: "Listen event commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	eventCmd.AddCommand(newListenTxCmd())
	eventCmd.AddCommand(newListenBlockCmd())
	eventCmd.AddCommand(newListenChaincodeCmd())
	return eventCmd
}

func newListenTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Listen to tx events",
		Run: func(c *cobra.Command, args []string) {
			//newConfig := cmd.NewConfig(cmd.Flags())
			txID, err := c.Flags().GetString(cmd.TxIDFlag)
			if err != nil || txID == "" {
				fmt.Println("[ txid can't empty ]", err)
				c.HelpFunc()(c, args)
				return
			}
			event, err := NewEventAction(api.ConfigFlags(c.Flags()))
			if err != nil {
				panic(err.Error())
			}
			if err = event.ListenTx(txID); err != nil {
				panic(err.Error())
			}
		},
	}
	cmd.InitTxID(txCmd.Flags())
	return txCmd
}

func newListenBlockCmd() *cobra.Command {
	blockCmd := &cobra.Command{
		Use:   "block",
		Short: "Listen to block events",
		Run: func(c *cobra.Command, args []string) {
			number, _ := c.Flags().GetUint64(cmd.BlockNumFlag)
			event, err := NewEventAction(api.ConfigFlags(c.Flags()), number)
			if err != nil {
				panic(err.Error())
			}
			if err = event.ListenBlock(); err != nil {
				panic(err.Error())
			}
		},
	}
	cmd.InitBlockNum(blockCmd.Flags())
	return blockCmd
}

func newListenChaincodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "chaincode",
		Short: "Listen to chaincode events",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}
