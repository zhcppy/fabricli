/*
@Time 2019-09-04 16:19
@Author ZH

*/
package actions

import (
	"math/rand"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	mspImpl "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/orderer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"github.com/zhcppy/fabricli/logger"
)

type Action struct {
	sdk                  *fabsdk.FabricSDK
	client               context.Client
	Peers                []fab.Peer
	Orderers             []fab.Orderer
	PeersByOrg           map[string]map[string]fab.Peer
	mspClientByOrg       map[string]*msp.Client
	clientProviderByUser map[string]context.ClientProvider
}

func New(configFile, provider string) (action *Action, err error) {
	action = &Action{
		PeersByOrg:           map[string]map[string]fab.Peer{},
		mspClientByOrg:       map[string]*msp.Client{},
		clientProviderByUser: map[string]context.ClientProvider{},
	}
	action.sdk, err = fabsdk.New(config.FromFile(configFile), newProviderOption(provider)...)
	if err != nil {
		return nil, errors.Errorf("Error initializing SDK: %s", err)
	}

	action.client, err = action.sdk.Context()()
	if err != nil {
		return nil, errors.WithMessage(err, "Error creating anonymous provider")
	}

	if err = action.OrgPeer(); err != nil {
		return
	}
	if err = action.FilterOrderers(); err != nil {
		return
	}
	logger.L().Debug("New fabsdk successfully")
	return action, nil
}

func (action *Action) OrgPeer() (err error) {
	organizationConfigs := action.client.EndpointConfig().NetworkConfig().Organizations

	for orgID, orgConfig := range organizationConfigs {
		logger.L().Debugf("- Org: %s, MSPID: %s", orgID, orgConfig.MSPID)
		action.PeersByOrg[orgID] = make(map[string]fab.Peer)
		for i, peerKey := range orgConfig.Peers {
			peerConfig, ok := action.client.EndpointConfig().PeerConfig(peerKey)
			if !ok {
				return errors.Errorf("failed to get peer configs for url [%s]", orgID)
			}
			peer, err := action.client.InfraProvider().CreatePeerFromConfig(&fab.NetworkPeer{PeerConfig: *peerConfig, MSPID: orgConfig.MSPID})
			if err != nil {
				return errors.Wrapf(err, "failed to create peer from config")
			}
			logger.L().Debugf("-- Peer[%d]: MSPID: %s, URL: %s, %s", i, peer.MSPID(), peer.URL(), peerKey)
			action.PeersByOrg[orgID][peerKey] = peer
			action.Peers = append(action.Peers, peer)
		}
	}
	return nil
}

func (action *Action) RandomPeer() (fab.Peer, error) {
	if len(action.Peers) == 0 {
		return nil, errors.New("No orders found")
	}
	return action.Peers[rand.Intn(len(action.Peers))], nil
}

func (action *Action) GetPeers() (peers map[string]fab.Peer) {
	peers = map[string]fab.Peer{}
	for _, v := range action.PeersByOrg {
		peers = v
	}
	return
}

func (action *Action) FilterOrderers(ordererURLs ...string) (err error) {
	ordererConfigs := action.client.EndpointConfig().OrderersConfig()

	if len(ordererConfigs) > 0 {
		logger.L().Debugf("- Orderer")
	}
	for i, ordererCfg := range ordererConfigs {
		logger.L().Debugf("-- Orderer[%d]: URL: %s %v", i, ordererCfg.URL, ordererCfg.GRPCOptions["ssl-target-name-override"])
		if len(ordererURLs) == 0 || ordererURLs[0] == ordererCfg.URL {
			newOrderer, err := orderer.New(action.client.EndpointConfig(), orderer.FromOrdererConfig(&ordererCfg))
			if err != nil {
				return errors.WithMessage(err, "creating orderer failed")
			}
			action.Orderers = append(action.Orderers, newOrderer)
		}
	}

	for name, ch := range action.client.EndpointConfig().NetworkConfig().Channels {
		logger.L().Debugf("- Channel: %s", name)
		for i, or := range ch.Orderers {
			logger.L().Debugf("-- Orderer[%d]: %s", i, or)
		}
		for peer := range ch.Peers {
			logger.L().Debugf("-- Peer: %s", peer)
		}
	}
	return nil
}

func (action *Action) RandomOrderer() (fab.Orderer, error) {
	if len(action.Orderers) == 0 {
		return nil, errors.New("No orders found")
	}
	return action.Orderers[rand.Intn(len(action.Orderers))], nil
}

func (action *Action) NewMspClient(orgID string) (*msp.Client, error) {
	if client, ok := action.mspClientByOrg[orgID]; ok {
		return client, nil
	}
	mspClient, err := msp.New(action.sdk.Context(), msp.WithOrg(orgID))
	if err != nil {
		return nil, errors.Errorf("error creating MSP client: %s", err)
	}
	action.mspClientByOrg[orgID] = mspClient
	return mspClient, nil
}

func (action *Action) ClientProvider(user mspImpl.SigningIdentity) (context.ClientProvider, error) {
	key := user.Identifier().MSPID + "_" + user.Identifier().ID
	if provider, ok := action.clientProviderByUser[key]; ok {
		return provider, nil
	}
	clientProvider := action.sdk.Context(fabsdk.WithIdentity(user))
	action.clientProviderByUser[key] = clientProvider
	return clientProvider, nil
}

// ====== User ====== //

//func (action *Action) NewUser(orgID, username, pwd string) (mspImpl.SigningIdentity, error) {
//	logger.L().Infof("Enrolling user %s...", username)
//
//	mspClient, err := action.NewMspClient(orgID)
//	if err != nil {
//		return nil, err
//	}
//
//	logger.L().Infof("Creating new user %s...", username)
//	err = mspClient.Enroll(username, msp.WithSecret(pwd))
//	if err != nil {
//		return nil, errors.Errorf("Enroll returned error: %v", err)
//	}
//
//	user, err := mspClient.GetSigningIdentity(username)
//	if err != nil {
//		return nil, errors.Errorf("GetSigningIdentity returned error: %v", err)
//	}
//
//	logger.L().Infof("Returning user [%s], MSPID [%s]", user.Identifier().ID, user.Identifier().MSPID)
//	return user, nil
//}

func (action *Action) User(orgID, peerUrl, username string) (mspImpl.SigningIdentity, error) {
	if orgID != "" && username != "" {
		return action.UserByOrg(orgID, username)
	}
	if peerUrl != "" && username != "" {
		return action.UserByPeer(peerUrl, username)
	}
	return nil, errors.Errorf("orgID and peerUrl cannot be empty at the same time")
}

func (action *Action) UserByPeer(peerUrl, username string) (mspImpl.SigningIdentity, error) {
	for orgID, peers := range action.PeersByOrg {
		for key, peer := range peers {
			if key == peerUrl || peer.URL() == peerUrl {
				return action.UserByOrg(orgID, username)
			}
		}
	}
	return nil, errors.Errorf("no found peer %s", peerUrl)
}

func (action *Action) UserByOrg(orgID, username string) (mspImpl.SigningIdentity, error) {
	if username == "" || orgID == "" {
		return nil, errors.Errorf("no username or orgID specified")
	}
	mspClient, err := action.NewMspClient(orgID)
	if err != nil {
		return nil, errors.WithMessage(err, "orgID: "+orgID)
	}

	user, err := mspClient.GetSigningIdentity(username)
	if err != nil {
		return nil, errors.Errorf("msp signing identity returned error: %s, username: %s", err.Error(), username)
	}

	logger.L().Infof("Returning user [%s], MSPID [%s]", user.Identifier().ID, user.Identifier().MSPID)
	return user, nil
}

func (action *Action) ChannelProvider(channelID string, user mspImpl.SigningIdentity) (context.ChannelProvider, error) {
	logger.L().Debugf("creating channel provider for user [%s] in org [%s]...", user.Identifier().ID, user.Identifier().MSPID)
	cp, err := action.ClientProvider(user)
	if err != nil {
		return nil, errors.Errorf("error getting session for user [%s,%s]: %v", user.Identifier().MSPID, user.Identifier().ID, err)
	}
	channelProvider := func() (context.Channel, error) {
		return contextImpl.NewChannel(cp, channelID)
	}
	return channelProvider, nil
}

// ====== Channel Client ====== //

func (action *Action) ChannelClient(channelID string, user mspImpl.SigningIdentity) (*channel.Client, error) {
	logger.L().Debugf("creating channel client for user [%s] in org [%s]...", user.Identifier().ID, user.Identifier().MSPID)
	provider, err := action.ChannelProvider(channelID, user)
	if err != nil {
		return nil, err
	}
	c, err := channel.New(provider)
	if err != nil {
		return nil, errors.Errorf("error creating new resmgmt client for user [%s,%s]: %v", user.Identifier().MSPID, user.Identifier().ID, err)
	}
	return c, nil
}

// ====== Event Client ====== //

func (action *Action) EventClient(channelID string, user mspImpl.SigningIdentity, opts ...event.ClientOption) (*event.Client, error) {
	logger.L().Debugf("creating event client for user [%s] in org [%s]...", user.Identifier().ID, user.Identifier().MSPID)
	channelProvider, err := action.ChannelProvider(channelID, user)
	if err != nil {
		return nil, errors.Errorf("error creating channel provider: %v", err)
	}
	c, err := event.New(channelProvider, opts...)
	if err != nil {
		return nil, errors.Errorf("error creating new event client: %v", err)
	}
	return c, nil
}

// ====== Event Client ====== //

func (action *Action) LedgerClient(channelID string, user mspImpl.SigningIdentity) (*ledger.Client, error) {
	logger.L().Debugf("creating ledger client for user [%s] in org [%s]...", user.Identifier().ID, user.Identifier().MSPID)
	channelProvider, err := action.ChannelProvider(channelID, user)
	if err != nil {
		return nil, errors.Errorf("error creating channel provider: %v", err)
	}
	c, err := ledger.New(channelProvider)
	if err != nil {
		return nil, errors.Errorf("error creating new ledger client: %v", err)
	}
	return c, nil
}

// ====== ResMgmt Client ====== //

func (action *Action) ResourceMgmtClient(user mspImpl.SigningIdentity) (*resmgmt.Client, error) {
	logger.L().Debugf("creating resmgmt client for user [%s] in org [%s]...", user.Identifier().ID, user.Identifier().MSPID)
	cp, err := action.ClientProvider(user)
	if err != nil {
		return nil, errors.Errorf("error getting session for user [%s,%s]: %v", user.Identifier().MSPID, user.Identifier().ID, err)
	}
	c, err := resmgmt.New(cp)
	if err != nil {
		return nil, errors.Errorf("error creating new resmgmt client for user [%s,%s]: %v", user.Identifier().MSPID, user.Identifier().ID, err)
	}
	return c, nil
}

// ====== Context Local ====== //

func (action *Action) LocalContext(user mspImpl.SigningIdentity) (context.Local, error) {
	logger.L().Debugf("creating local for user [%s] in org [%s]...", user.Identifier().ID, user.Identifier().MSPID)
	cp, err := action.ClientProvider(user)
	if err != nil {
		return nil, errors.Errorf("error getting context for user [%s,%s]: %v", user.Identifier().MSPID, user.Identifier().ID, err)
	}
	local, err := contextImpl.NewLocal(cp)
	if err != nil {
		return nil, errors.Errorf("error creating new local for user [%s,%s]: %v", user.Identifier().MSPID, user.Identifier().ID, err)
	}
	return local, nil
}

func (action *Action) Close() error {
	if action.sdk != nil {
		action.sdk.Close()
	}
	return nil
}
