package channel

import (
	"fmt"

	"github.com/zhcppy/fabricli/api"

	"github.com/spf13/cobra"
)

const (
	create = "create"
	join   = "join"
)

func NewCmd() *cobra.Command {
	var channelCmd = &cobra.Command{
		Use:   "channel",
		Short: "Channel commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	channelCmd.AddCommand(&cobra.Command{Use: create, Run: channelRun})
	channelCmd.AddCommand(&cobra.Command{Use: join, Run: channelRun})
	return channelCmd
}

func channelRun(cmd *cobra.Command, args []string) {
	config := api.GetConfig()
	channel, err := NewChannelAction(config)
	if err != nil {
		panic(err)
	}
	if cmd.Use == create {
		tx, err := channel.Create()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Success to create channel [%s], TX [%s]", channel.ChannelID, tx)
	}
	if cmd.Use == join {
		err := channel.Join()
		if err != nil {
			panic(err)
		}
	}
}
