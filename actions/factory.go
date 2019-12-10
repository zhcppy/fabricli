/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package actions

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/common/selection/dynamicselection"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/common/selection/fabricselection"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/common/selection/staticselection"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/options"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite/bccsp/multisuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/factory/defcore"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/factory/defsvc"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/provider/chpvdr"
	"github.com/pkg/errors"
	"github.com/zhcppy/fabricli/logger"
)

const (
	// AutoDetectSelectionProvider indicates that a selection provider is to be automatically determined using channel capabilities
	AutoDetectSelectionProvider = "auto"

	// StaticSelectionProvider indicates that a static selection provider is to be used for selecting Peers for invoke/query commands
	StaticSelectionProvider = "static"

	// DynamicSelectionProvider indicates that a dynamic selection provider is to be used for selecting Peers for invoke/query commands
	DynamicSelectionProvider = "dynamic"

	// FabricSelectionProvider indicates that the Fabric selection provider is to be used for selecting Peers for invoke/query commands
	FabricSelectionProvider = "fabric"
)

func newProviderOption(selectionProvider string) (opts []fabsdk.Option) {
	if selectionProvider != "" && selectionProvider != AutoDetectSelectionProvider {
		svcPackage, _ := newServiceProviderFactory(selectionProvider)
		opts = append(opts, fabsdk.WithServicePkg(svcPackage))
	}
	opts = append(opts, fabsdk.WithCorePkg(newcryptoSuiteProvider()))
	return
}

// serviceProviderFactory is configured with either static or dynamic selection provider
type serviceProviderFactory struct {
	defsvc.ProviderFactory
	selectionProvider string
}

func newServiceProviderFactory(selectionProvider string) (*serviceProviderFactory, error) {
	return &serviceProviderFactory{selectionProvider: selectionProvider}, nil
}

type fabricSelectionChannelProvider struct {
	fab.ChannelProvider
	service           fab.ChannelService
	selection         fab.SelectionService
	selectionProvider string
}

type fabricSelectionChannelService struct {
	fab.ChannelService
	selection fab.SelectionService
}

// CreateChannelProvider returns a new default implementation of channel provider
func (f *serviceProviderFactory) CreateChannelProvider(config fab.EndpointConfig, opts ...options.Opt) (fab.ChannelProvider, error) {
	chProvider, err := chpvdr.New(config, opts...)
	if err != nil {
		return nil, err
	}
	return &fabricSelectionChannelProvider{
		ChannelProvider:   chProvider,
		selectionProvider: f.selectionProvider,
	}, nil
}

type closable interface {
	Close()
}

// Close frees resources and caches.
func (cp *fabricSelectionChannelProvider) Close() {
	if c, ok := cp.ChannelProvider.(closable); ok {
		c.Close()
	}
	if cp.selection != nil {
		if c, ok := cp.selection.(closable); ok {
			c.Close()
		}
	}
}

type providerInit interface {
	Initialize(providers context.Providers) error
}

func (cp *fabricSelectionChannelProvider) Initialize(providers context.Providers) error {
	if init, ok := cp.ChannelProvider.(providerInit); ok {
		return init.Initialize(providers)
	}
	return nil
}

// ChannelService creates a ChannelService for an identity
func (cp *fabricSelectionChannelProvider) ChannelService(ctx fab.ClientContext, channelID string) (fab.ChannelService, error) {
	chService, err := cp.ChannelProvider.ChannelService(ctx, channelID)
	if err != nil {
		return nil, err
	}

	discovery, err := chService.Discovery()
	if err != nil {
		return nil, err
	}

	if cp.selection == nil {
		switch cp.selectionProvider {
		case StaticSelectionProvider:
			logger.L().Debug("Using static selection provider.")
			cp.selection, err = staticselection.NewService(discovery)
		case DynamicSelectionProvider:
			logger.L().Debug("Using dynamic selection provider.")
			cp.selection, err = dynamicselection.NewService(ctx, channelID, discovery)
		case FabricSelectionProvider:
			logger.L().Debug("Using Fabric selection provider.")
			cp.selection, err = fabricselection.New(ctx, channelID, discovery)
		default:
			return nil, errors.Errorf("invalid selection provider: %s", cp.selectionProvider)
		}

		if err != nil {
			return nil, err
		}
	}

	return &fabricSelectionChannelService{
		ChannelService: chService,
		selection:      cp.selection,
	}, nil
}

func (cs *fabricSelectionChannelService) Selection() (fab.SelectionService, error) {
	return cs.selection, nil
}

// cryptoSuiteProviderFactory will provide custom cryptosuite (bccsp.BCCSP)
type cryptoSuiteProviderFactory struct {
	defcore.ProviderFactory
}

func newcryptoSuiteProvider() *cryptoSuiteProviderFactory {
	return &cryptoSuiteProviderFactory{}
}

// CreateCryptoSuiteProvider returns a new default implementation of BCCSP
func (f *cryptoSuiteProviderFactory) CreateCryptoSuiteProvider(config core.CryptoSuiteConfig) (core.CryptoSuite, error) {
	return multisuite.GetSuiteByConfig(config)
}
