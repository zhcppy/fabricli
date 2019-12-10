/*
@Time 2019-08-29 11:10
@Author ZH

*/
package api

import (
	"fmt"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zhcppy/fabricli/logger"
)

// Config Precedence: [explicit call to Set] > [flag] > [env] > [config] > [key/value store] > [default]

const (
	ProjectHome = "GABRICLI_HOME"
	DefaultHomePath = "$HOME/.fabricli"

	LoggerLevelTag          = "LogLevel"
	ConfigFileTag           = "ConfigFile"
	UserTag                 = "Username"
	ChannelIDTag            = "ChannelId"
	OrgIdTag                = "OrgId"
	PeerUrlTag              = "PeerUrl"
	OrdererUrlTag           = "OrdererUrl"
	SelectionProviderTag    = "SelectionProvider"
	CollectionConfigFileTag = "CollectionConfigFile"
	ChaincodeArgsTag        = "ChaincodeArgs"
	ChaincodeIDTag          = "ChaincodeId"
	ChaincodePathTag        = "ChaincodePath"
	ChaincodeVersionTag     = "ChaincodeVersion"
	ChaincodeEventTag       = "ChaincodeEvent"
	ChaincodePolicyTag      = "ChaincodePolicy"
	GoPathTag               = "GoPath"
)

func init() {
	_ = os.Setenv(ProjectHome, "github.com/zhcppy/fabricli")
	viper.SetConfigName("fabricli")
	viper.AddConfigPath(DefaultHomePath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$GOPATH/src/github.com/zhcppy/fabricli")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.L().Warn("Config file not found;")
		} else {
			logger.L().Warn("Config file was found but another error was produced")
		}
		return
	}
	logger.L().Debug("Config file found and successfully parsed", viper.ConfigFileUsed())
	viper.WatchConfig()
	viper.OnConfigChange(defConfig.onchange)
}

// default
var defConfig = &Config{
	SelectionProvider: "auto",
	CCodeInfo: CCodeInfo{
		ChaincodeVersion: "v0",
	},
}
var once = &sync.Once{}

type Config struct {
	SelectionProvider string    `json:"SelectionProvider"`
	ConfigFile        string    `json:"ConfigFile"`
	Username          string    `json:"Username"`
	OrgID             string    `json:"OrgId"`
	PeerUrl           string    `json:"PeerUrl"`
	OrdererURL        string    `json:"OrdererUrl"`
	ChannelID         string    `json:"ChannelId"`
	CCodeInfo         CCodeInfo `json:"CCodeInfo"`
}

type CCodeInfo struct {
	CollectionConfigFile string `json:"CollectionConfigFile"`
	ChaincodeArgs        string `json:"ChaincodeArgs"`
	ChaincodeID          string `json:"ChaincodeId"`
	ChaincodePath        string `json:"ChaincodePath"`
	ChaincodeVersion     string `json:"ChaincodeVersion"`
	ChaincodeEvent       string `json:"ChaincodeEvent"`
	ChaincodePolicy      string `json:"ChaincodePolicy"`
	GoPath               string `json:"GoPath"`
}

func (c *Config) check() *Config {
	var err error
	switch {
	case c.Username == "":
		err = fmt.Errorf("UserName can't empty")
	case c.ChannelID == "":
		err = fmt.Errorf("ChannelID can't empty")
	case c.SelectionProvider == "":
		err = fmt.Errorf("SelectionProvider can't empty")
	case c.ConfigFile == "":
		err = fmt.Errorf("ConfigFile can't empty")
	default:
		return c
	}
	panic(err.Error())
}

func (c *Config) onchange(e fsnotify.Event) {
	fmt.Println("Config file changed:", e.String())
}

func GetConfig() *Config {
	once.Do(viperCfg)
	return defConfig.check()
}

func ConfigFlags(flags *pflag.FlagSet) *Config {
	once.Do(func() {
		flags.Visit(func(flag *pflag.Flag) {
			viper.BindPFlag(flag.Name, flag)
		})
		viperCfg()
	})
	return defConfig.check()
}

func viperCfg() {
	logger.SetLevel(viper.GetInt(LoggerLevelTag))
	if err := viper.Unmarshal(defConfig); err != nil {
		fmt.Println("Failed to unmarshal config file", err.Error())
	}
	//viper.Debug()
	//logger.L().Debug(defConfig)
}
