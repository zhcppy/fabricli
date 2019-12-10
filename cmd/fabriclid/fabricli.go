/*
@Time 2019-08-29 11:36
@Author ZH

*/
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zhcppy/fabricli/api/chaincode"
	"github.com/zhcppy/fabricli/api/channel"
	"github.com/zhcppy/fabricli/api/event"
	"github.com/zhcppy/fabricli/api/query"
	"github.com/zhcppy/fabricli/cmd"
	"github.com/zhcppy/fabricli/console"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "fabricli",
		Short: "Fabric Client",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	flags := rootCmd.PersistentFlags()
	cmd.InitLoggingLevel(flags)
	cmd.InitConfigFile(flags)
	cmd.InitUserName(flags)
	cmd.InitOrgIDs(flags)
	cmd.InitPeerURL(flags)
	cmd.InitOrdererURL(flags)
	cmd.InitChannelID(flags)
	cmd.InitSelectionProvider(flags)
	cmd.InitOrdererTLSCertificate(flags)

	rootCmd.AddCommand(console.NewCmd())
	rootCmd.AddCommand(query.NewCmd())
	rootCmd.AddCommand(event.NewCmd())
	rootCmd.AddCommand(channel.NewCmd())
	rootCmd.AddCommand(chaincode.NewCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
