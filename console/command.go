/*
@Time 2019-08-29 11:36
@Author ZH

*/
package console

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	consoleCmd := &cobra.Command{
		Use:   "console",
		Short: "Fabric command-line console implemented by Golang",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}
	return consoleCmd
}
