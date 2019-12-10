/*
@Time 2019-08-29 16:36
@Author ZH

*/
package channel

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/pkg/errors"
	"github.com/zhcppy/fabricli/actions"
	"github.com/zhcppy/fabricli/api"
	"github.com/zhcppy/fabricli/logger"
)

type Channel struct {
	OrgID     string
	ChannelID string

	user    msp.SigningIdentity
	orderer fab.Orderer
	client  *resmgmt.Client
	action  *actions.Action
}

func NewChannelAction(c *api.Config) (*Channel, error) {
	action, err := actions.New(c.ConfigFile, c.SelectionProvider)
	if err != nil {
		return nil, err
	}
	user, err := action.User(c.OrgID, c.PeerUrl, c.Username)
	if err != nil {
		return nil, err
	}

	client, err := action.ResourceMgmtClient(user)
	if err != nil {
		return nil, err
	}

	orderer, err := action.RandomOrderer()
	if err != nil {
		return nil, err
	}

	return &Channel{
		action:  action,
		client:  client,
		orderer: orderer,
	}, nil
}

func (c *Channel) Create() (txID fab.TransactionID, err error) {
	logger.L().Infof("Attempting to create/update channel: %s", c.ChannelID)

	req := resmgmt.SaveChannelRequest{ChannelID: c.ChannelID, SigningIdentities: []msp.SigningIdentity{c.user}}
	resp, err := c.client.SaveChannel(req, resmgmt.WithOrderer(c.orderer))
	if err != nil {
		return "", errors.Errorf("Error from save channel: %s", err.Error())
	}
	return resp.TransactionID, nil
}

func (c *Channel) Join() error {
	logger.L().Debugf("Attempting to join channel: %s\n", c.ChannelID)

	var peers []fab.Peer
	for _, peer := range c.action.PeersByOrg[c.OrgID] {
		logger.L().Debugf("Joining channel [%s] on org [%s] peers:\n", c.ChannelID, c.OrgID)
		for _, peer := range peers {
			fmt.Printf("-- %s\n", peer.URL())
		}
		peers = append(peers, peer)
	}

	logger.L().Debugf("==========> JOIN ORG: %s\n", c.OrgID)
	if len(peers) == 0 {
		return errors.Errorf("at least one peer is required for join")
	}

	if err := c.client.JoinChannel(c.ChannelID, resmgmt.WithTargets(peers...), resmgmt.WithOrderer(c.orderer)); err != nil {
		return errors.WithMessage(err, "Could not join channel: %v")
	}

	logger.L().Debugf("Channel %s joined!\n", c.ChannelID)
	return nil
}
