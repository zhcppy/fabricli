/*
@Time 2019-08-29 11:36
@Author ZH

*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zhcppy/fabricli/api"
)

// InitLoggingLevel initializes the logging level from the provided arguments
func InitLoggingLevel(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		loggingLevelFlag        = "level"
		loggingLevelDescription = "Logging level [-1 debug, 0 info, 1 warn, 2 error, 3 D-panic, 4 panic, 5 fatal]"
		defaultLoggingLevel     = "0"
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultLoggingLevel, loggingLevelDescription, defaultValueAndDescription...)
	value, err := strconv.Atoi(defaultValue)
	if err != nil {
		fmt.Printf("Invalid number for [%s]: %s\n", loggingLevelFlag, defaultValue)
	}
	flags.Int(loggingLevelFlag, value, description)
	viper.RegisterAlias(loggingLevelFlag, api.LoggerLevelTag)
	//viper.BindPFlag(api.LoggingLevelTag, flags.Lookup(loggingLevelFlag))
}

// InitConfigFile initializes the config file path from the provided arguments
func InitConfigFile(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		configFileFlag        = "config"
		configFileDescription = "The path of the config.yaml file"
		defaultConfigFile     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultConfigFile, configFileDescription, defaultValueAndDescription...)
	flags.String(configFileFlag, defaultValue, description)
	viper.RegisterAlias(configFileFlag, api.ConfigFileTag)
	//viper.BindPFlag(api.ConfigFileTag, flags.Lookup(configFileFlag))
}

// InitUserName initializes the user name from the provided arguments
func InitUserName(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		userFlag        = "user"
		userDescription = "The user name"
		defaultUser     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultUser, userDescription, defaultValueAndDescription...)
	flags.String(userFlag, defaultValue, description)
	viper.RegisterAlias(userFlag, api.UserTag)
	//viper.BindPFlag(api.UserTag, flags.Lookup(userFlag))
}

// InitChannelID initializes the channel ID from the provided arguments
func InitChannelID(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		channelIDFlag        = "cid"
		channelIDDescription = "The channel ID"
		defaultChannelID     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultChannelID, channelIDDescription, defaultValueAndDescription...)
	flags.String(channelIDFlag, defaultValue, description)
	viper.RegisterAlias(channelIDFlag, api.ChannelIDTag)
	//viper.BindPFlag(api.ChannelIDTag, flags.Lookup(channelIDFlag))
}

// InitOrgIDs initializes the org IDs from the provided arguments
func InitOrgIDs(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		orgIDsFlag        = "orgid"
		orgIDsDescription = "A comma-separated list of organization IDs"
		defaultOrgIDs     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultOrgIDs, orgIDsDescription, defaultValueAndDescription...)
	flags.String(orgIDsFlag, defaultValue, description)
	//viper.BindPFlag(api.OrgIdTag, flags.Lookup(orgIDsFlag))
}

// InitPeerURL initializes the peer URL from the provided arguments
func InitPeerURL(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		peerURLFlag        = "peer"
		peerURLDescription = "A comma-separated list of peer targets, e.g. 'grpcs://localhost:7051,grpcs://localhost:8051'"
		defaultPeerURL     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultPeerURL, peerURLDescription, defaultValueAndDescription...)
	flags.String(peerURLFlag, defaultValue, description)
	viper.RegisterAlias(peerURLFlag, api.PeerUrlTag)
	//viper.BindPFlag(api.PeerUrlTag, flags.Lookup(peerURLFlag))
}

// InitOrdererURL initializes the orderer URL from the provided arguments
func InitOrdererURL(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		ordererFlag           = "orderer"
		ordererURLDescription = "The URL of the orderer, e.g. grpcs://localhost:7050"
		defaultOrdererURL     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultOrdererURL, ordererURLDescription, defaultValueAndDescription...)
	flags.String(ordererFlag, defaultValue, description)
	viper.RegisterAlias(ordererFlag, api.OrdererUrlTag)
	//viper.BindPFlag(api.OrdererUrlTag, flags.Lookup(ordererFlag))
}

// InitChaincodeID initializes the chaincode ID from the provided arguments
func InitChaincodeID(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		chaincodeIDFlag        = "ccid"
		chaincodeIDDescription = "The Chaincode ID"
		defaultChaincodeID     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultChaincodeID, chaincodeIDDescription, defaultValueAndDescription...)
	flags.String(chaincodeIDFlag, defaultValue, description)
	viper.RegisterAlias(chaincodeIDFlag, api.ChaincodeIDTag)
	//viper.BindPFlag(api.ChaincodeIDTag, flags.Lookup(chaincodeIDFlag))
}

// InitChaincodeEvent initializes the chaincode event name from the provided arguments
func InitChaincodeEvent(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		chaincodeEventFlag        = "event"
		chaincodeEventDescription = "The name of the chaincode event to listen for"
		defaultChaincodeEvent     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultChaincodeEvent, chaincodeEventDescription, defaultValueAndDescription...)
	flags.String(chaincodeEventFlag, defaultValue, description)
	viper.RegisterAlias(chaincodeEventFlag, api.ChaincodeEventTag)
	//viper.BindPFlag(api.ChaincodeEventTag, flags.Lookup(chaincodeEventFlag))
}

// InitChaincodePath initializes the chaincode install source path from the provided arguments
func InitChaincodePath(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		chaincodePathFlag        = "ccp"
		chaincodePathDescription = "The Chaincode path"
		defaultChaincodePath     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultChaincodePath, chaincodePathDescription, defaultValueAndDescription...)
	flags.String(chaincodePathFlag, defaultValue, description)
	viper.RegisterAlias(chaincodePathFlag, api.ChaincodePathTag)
	//viper.BindPFlag(api.ChaincodePathTag, flags.Lookup(chaincodePathFlag))
}

// InitChaincodePolicy initializes the chaincode policy from the provided arguments
func InitChaincodePolicy(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		chaincodePolicyFlag        = "policy"
		chaincodePolicyDescription = "The chaincode policy, e.g. OutOf(1,'Org1MSP.admin','Org2MSP.admin',AND('Org3MSP.member','Org4MSP.member'))"
		defaultChaincodePolicy     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultChaincodePolicy, chaincodePolicyDescription, defaultValueAndDescription...)
	flags.String(chaincodePolicyFlag, defaultValue, description)
	viper.RegisterAlias(chaincodePolicyFlag, api.ChaincodePolicyTag)
	//viper.BindPFlag(api.ChaincodePolicyTag, flags.Lookup(chaincodePolicyFlag))
}

// InitChaincodeVersion initializes the chaincode version from the provided arguments
func InitChaincodeVersion(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		chaincodeVersionFlag        = "v"
		chaincodeVersionDescription = "The chaincode version"
		defaultChaincodeVersion     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultChaincodeVersion, chaincodeVersionDescription, defaultValueAndDescription...)
	flags.String(chaincodeVersionFlag, defaultValue, description)
	viper.RegisterAlias(chaincodeVersionFlag, api.ChaincodeVersionTag)
	//viper.BindPFlag(api.ChaincodeVersionTag, flags.Lookup(chaincodeVersionFlag))
}

// InitChaincodeArgs initializes the invoke/query args from the provided arguments
func InitChaincodeArgs(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	//Note that $rand(N) may be used anywhere within the value of the arg in order to generate a random value between 0 and N. For example {"Func":"function","Args":["arg_$rand(100)","$rand(10)"]}.
	const (
		chaincodeArgsFlag = "args"
		argsDescription   = `The args in JSON format. Example: {"Func":"function","Args":["arg1","arg2"]}.`
		defaultArgsFlag   = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultArgsFlag, argsDescription, defaultValueAndDescription...)
	flags.String(chaincodeArgsFlag, defaultValue, description)
	viper.RegisterAlias(chaincodeArgsFlag, api.ChaincodeArgsTag)
	//viper.BindPFlag(api.ChaincodeArgsTag, flags.Lookup(chaincodeArgsFlag))
}

// InitTxFile initializes the path of the .tx file used to create/update a channel from the provided arguments
func InitTxFile(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		txFileFlag        = "txfile"
		txFileDescription = "The path of the channel.tx file"
		defaultTxFile     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultTxFile, txFileDescription, defaultValueAndDescription...)
	flags.String(txFileFlag, defaultValue, description)
	//viper.BindPFlag(api.tx, flags.Lookup(txFileFlag))
}

const TxIDFlag = "txid"

// InitTxID initializes the transaction D from the provided arguments
func InitTxID(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		txIDDescription = "The transaction ID"
		defaultTxID     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultTxID, txIDDescription, defaultValueAndDescription...)
	flags.String(TxIDFlag, defaultValue, description)
}

const BlockNumFlag = "number"

// InitBlockNum initializes the bluck number from the provided arguments
func InitBlockNum(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		blockNumDescription = "The block number"
		defaultBlockNum     = "0"
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultBlockNum, blockNumDescription, defaultValueAndDescription...)
	value, err := strconv.ParseUint(defaultValue, 10, 64)
	if err != nil {
		fmt.Printf("Invalid number for [%s]: %s\n", BlockNumFlag, defaultValue)
	}
	flags.Uint64(BlockNumFlag, value, description)
}

const BlockHashFlag = "hash"

// InitBlockHash initializes the block hash from the provided arguments
func InitBlockHash(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		blockHashDescription = "The block hash"
		defaultBlockHash     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultBlockHash, blockHashDescription, defaultValueAndDescription...)
	flags.String(BlockHashFlag, defaultValue, description)
}

// InitCollectionConfigFile initializes the collection config file from the provided arguments
func InitCollectionConfigFile(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		collectionConfigFileFlag        = "collconfig"
		collectionConfigFileDescription = "The path of the JSON file that contains the private data collection configuration for the chaincode"
		defaultCollectionConfigFile     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultCollectionConfigFile, collectionConfigFileDescription, defaultValueAndDescription...)
	flags.String(collectionConfigFileFlag, defaultValue, description)
	viper.RegisterAlias(collectionConfigFileFlag, api.CollectionConfigFileTag)
	//viper.BindPFlag(api.CollectionConfigFileTag, flags.Lookup(collectionConfigFileFlag))
}

// InitSelectionProvider initializes the peer selection provider from the provided arguments
func InitSelectionProvider(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	//The possible values are: (1) static - Selects all peers; (2) dynamic - Uses the built-in selection service from the SDK to select a minimal set of peers according to the endorsement policy of the chaincode; (3) fabric - Uses Fabric's Discovery Service to select a minimal set of peers according to the endorsement/collection policy of the chaincode; (4) auto (default) - Automatically determines which selection service to use based on channel capabilities.
	const (
		selectionProviderFlag        = "provider"
		selectionProviderDescription = "The peer selection provider for invoke/query commands. [ static, dynamic, fabric, auto(default) ]"
		defaultSelectionProvider     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultSelectionProvider, selectionProviderDescription, defaultValueAndDescription...)
	flags.String(selectionProviderFlag, defaultValue, description)
	viper.RegisterAlias(selectionProviderFlag, api.SelectionProviderTag)
	//viper.BindPFlag(api.SelectionProviderTag, flags.Lookup(selectionProviderFlag))
}

// InitOrdererTLSCertificate initializes the orderer TLS certificate from the provided arguments
func InitOrdererTLSCertificate(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		certificateFileFlag    = "cacert"
		certificateDescription = "The path of the ca-cert.pem file"
		defaultCertificate     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultCertificate, certificateDescription, defaultValueAndDescription...)
	flags.String(certificateFileFlag, defaultValue, description)
	//viper.BindPFlag(certificateFileFlag, flags.Lookup(certificateFileFlag))
}

// InitGoPath initializes the gopath from the provided arguments
func InitGoPath(flags *pflag.FlagSet, defaultValueAndDescription ...string) {
	const (
		goPathFlag        = "gopath"
		goPathDescription = "GOPATH for chaincode install command. If not set, GOPATH is taken from the environment"
		defaultGoPath     = ""
	)
	defaultValue, description := GetDefaultValueAndDescription(defaultGoPath, goPathDescription, defaultValueAndDescription...)
	flags.String(goPathFlag, defaultValue, description)
	//viper.BindPFlag(api.GoPathTag, flags.Lookup(goPathFlag))
}

func GetDefaultValueAndDescription(defaultValue string, defaultDescription string, overrides ...string) (value, description string) {
	if len(overrides) > 0 {
		value = overrides[0]
	} else {
		value = defaultValue
	}
	if len(overrides) > 1 {
		description = overrides[1]
	} else {
		description = defaultDescription
	}
	return value, description
}
