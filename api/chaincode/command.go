package chaincode

import (
	"github.com/spf13/cobra"
	"github.com/zhcppy/fabricli/api"
	"github.com/zhcppy/fabricli/cmd"
)

func NewCmd() *cobra.Command {
	var chaincodeCmd = &cobra.Command{
		Use:   "chaincode",
		Short: "Chaincode commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}
	flags := chaincodeCmd.PersistentFlags()
	cmd.InitChaincodeArgs(flags)
	cmd.InitCollectionConfigFile(flags)
	cmd.InitChaincodeEvent(flags)
	cmd.InitChaincodeID(flags)
	cmd.InitChaincodePath(flags)
	cmd.InitChaincodePolicy(flags)
	cmd.InitChaincodeVersion(flags)
	cmd.InitGoPath(flags)

	chaincodeCmd.AddCommand(newCCInstallCmd())
	chaincodeCmd.AddCommand(newCCInfoCmd())
	chaincodeCmd.AddCommand(newCCInstantiateCmd())
	chaincodeCmd.AddCommand(newCCUpgradeCmd())
	chaincodeCmd.AddCommand(newCCInvokeCmd())
	return chaincodeCmd
}

func newCCInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install chaincode.",
		Run: func(cmd *cobra.Command, args []string) {
			//action, err := NewCCAction(api.GetConfig())
			//if err != nil {
			//	panic(err)
			//}
			//action.Install(config.CCodeInfo)

		},
	}
}

func newCCUpgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade chaincode.",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}

func newCCInstantiateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "instantiate",
		Short: "Instantiate chaincode.",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := api.GetConfig()
			action, err := NewCCAction(cfg)
			if err != nil {
				panic(err)
			}
			err = action.Instantiate(cfg.CCodeInfo)
			if err != nil {
				panic(err)
			}
		},
	}
}

func newCCInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Get chaincode info,Retrieves details about the chaincode",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := api.GetConfig()
			action, err := NewCCAction(cfg)
			if err != nil {
				panic(err)
			}
			err = action.QueryInfo(cfg.CCodeInfo.ChaincodeID, cfg.CCodeInfo.ChaincodeArgs)
			if err != nil {
				panic(err)
			}
		},
	}
}

func newCCInvokeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "invoke",
		Short: "invoke chaincode.",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := api.GetConfig()
			action, err := NewCCAction(cfg)
			if err != nil {
				panic(err)
			}
			err = action.Invoke(cfg.CCodeInfo.ChaincodeID, cfg.CCodeInfo.ChaincodeArgs)
			if err != nil {
				panic(err)
			}
		},
	}
}
