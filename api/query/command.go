package query

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhcppy/fabricli/api"
	"github.com/zhcppy/fabricli/console"
)

func NewCmd() *cobra.Command {
	var queryCmd = &cobra.Command{
		Use:       "query",
		Short:     "Query commands",
		Example:   "query info",
		ValidArgs: []string{"info"},
		Run: func(cmd *cobra.Command, args []string) {
			config := api.ConfigFlags(cmd.Flags())
			c, err := console.New(consoler{Config: config}, console.WithPrompt("> QueryAction."))
			if err != nil {
				fmt.Println("console error:", err.Error())
				return
			}
			for i, arg := range args {
				arg = fmt.Sprintf("%s%s()", string(arg[0]-32), arg[1:])
				fmt.Println("> QueryAction." + arg)
				if err := c.Execute(arg); err != nil {
					fmt.Printf("/ninput:[ %s ], execute err:%s", arg, err.Error())
				}
				if i+1 == len(args) {
					return
				}
			}
			c.Interactive()
			defer c.Close()
		},
	}
	return queryCmd
}
