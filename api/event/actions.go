package event

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/pkg/errors"
	"github.com/zhcppy/fabricli/actions"
	"github.com/zhcppy/fabricli/api"
	"github.com/zhcppy/fabricli/logger"
	"github.com/zhcppy/fabricli/printer"
)

type Event struct {
	action      *actions.Action
	eventClient *event.Client
	inputEvent
}

func NewEventAction(c *api.Config, numbers ...uint64) (*Event, error) {
	action, err := actions.New(c.ConfigFile, c.SelectionProvider)
	if err != nil {
		return nil, err
	}
	user, err := action.User(c.OrgID, c.PeerUrl, c.Username)
	if err != nil {
		return nil, err
	}
	logger.L().Debugf("new event action, user:%v", user.Identifier())
	var startNumber uint64
	if len(numbers) > 0 {
		startNumber = numbers[0]
		fmt.Printf("Listen block number start [%d]\n", startNumber)
	}

	eventClient, err := action.EventClient(c.ChannelID, user,
		event.WithBlockEvents(), event.WithBlockNum(startNumber), event.WithSeekType("from"))
	if err != nil {
		return nil, err
	}
	return &Event{
		action:      action,
		eventClient: eventClient,
		inputEvent:  newInputEvent(),
	}, nil
}

func (e *Event) ListenBlock() error {
	fmt.Println("Registering listen block event ...")

	blockReg, eventCh, err := e.eventClient.RegisterBlockEvent()
	if err != nil {
		return errors.WithMessage(err, "Error registering for block events")
	}
	defer e.eventClient.Unregister(blockReg)

	exit := e.WaitForEnter()
	for {
		select {
		case <-exit:
			fmt.Println("Listen block event exiting ...")
			return nil
		case block, ok := <-eventCh:
			if !ok {
				return errors.WithMessage(err, "unexpected closed channel while waiting for block event")
			}
			//fmt.Println(block)
			printer.Success("URL:%s, BlockNumber:%d \n", block.SourceURL, block.Block.Header.Number)
			//fmt.Println("Press <enter> to terminate")
		}
	}
}

func (e *Event) ListenTx(txID string) error {
	fmt.Printf("Registering listen TX event for TxID [%s]\n", txID)

	registration, eventCh, err := e.eventClient.RegisterTxStatusEvent(txID)
	if err != nil {
		return errors.WithMessage(err, "Error registering for block events")
	}
	defer e.eventClient.Unregister(registration)

	exit := e.WaitForEnter()
	select {
	case <-exit:
		fmt.Println("Listen TX event exiting ...")
		return nil
	case tx, ok := <-eventCh:
		if !ok {
			return errors.WithMessage(err, "unexpected closed channel while waiting for tx status event")
		}
		fmt.Println(tx)
		fmt.Printf("Received TX event. TxID: %s, Code: %s, Error: %s\n", tx.TxID, tx.TxValidationCode, err)
	}
	return nil
}

func (e *Event) ListenChaincode(chaincodeID, eventFilter string) error {
	registration, eventCh, err := e.eventClient.RegisterChaincodeEvent(chaincodeID, eventFilter)
	if err != nil {
		return errors.WithMessage(err, "Error registering for block events")
	}
	defer e.eventClient.Unregister(registration)

	exit := e.WaitForEnter()
	for {
		select {
		case <-exit:
			return nil
		case cc, ok := <-eventCh:
			if !ok {
				return errors.WithMessage(err, "unexpected closed channel while waiting for block event")
			}
			fmt.Println(cc)
			fmt.Println("Press <enter> to terminate")
		}
	}
}

func (e *Event) ListenFilteredBlock() error {
	fmt.Printf("Registering filtered block event\n")

	registration, eventCh, err := e.eventClient.RegisterFilteredBlockEvent()
	if err != nil {
		return errors.WithMessage(err, "Error registering for filtered block events")
	}
	defer e.eventClient.Unregister(registration)

	exit := e.WaitForEnter()
	for {
		select {
		case <-exit:
			return nil
		case block, ok := <-eventCh:
			if !ok {
				return errors.WithMessage(err, "unexpected closed channel while waiting for filtered block event")
			}
			fmt.Println(block)
			fmt.Println("Press <enter> to terminate")
		}
	}
}
