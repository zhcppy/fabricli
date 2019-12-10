package task

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/zhcppy/fabricli/logger"
)

type Task interface {
	// Invoke invokes the task
	Invoke()

	// Attempts returns the number of invocation attempts that were made
	// in order to achieve a successful response
	Attempts() int

	// LastError returns the last error that occurred
	LastError() error
}

type ChaincodeTask struct {
	ctxt          Context
	id            string
	channelClient *channel.Client
	targets       []fab.Peer
	chaincodeID   string
	args          ArgStruct
	retryOpts     retry.Opts
	startedCB     func()
	completedCB   func(err error)
	attempt       int
	lastErr       error
	isQuery       bool
}

func NewCCTask(ctxt Context, id string, channelClient *channel.Client, targets []fab.Peer, chaincodeID string,
	args ArgStruct, retryOpts retry.Opts, startedCB func(), completedCB func(err error), isQuery bool) *ChaincodeTask {
	return &ChaincodeTask{
		ctxt:          ctxt,
		id:            id,
		channelClient: channelClient,
		targets:       targets,
		chaincodeID:   chaincodeID,
		args:          args,
		retryOpts:     retryOpts,
		startedCB:     startedCB,
		completedCB:   completedCB,
		attempt:       1,
		isQuery:       isQuery,
	}
}

func (t *ChaincodeTask) Invoke() {
	if t.isQuery {
		t.Query()
	} else {
		t.doInvoke()
	}
}

func (t *ChaincodeTask) doInvoke() {
	t.startedCB()
	logger.L().Debugf("(%s) - Invoking chaincode: %s, function: %s, args: %+v. Attempt #%d...\n",
		t.id, t.chaincodeID, t.args.Func, t.args.Args, t.attempt)

	var opts []channel.RequestOption
	opts = append(opts, channel.WithRetry(t.retryOpts))
	opts = append(opts, channel.WithBeforeRetry(func(err error) {
		t.attempt++
	}))
	if len(t.targets) > 0 {
		opts = append(opts, channel.WithTargets(t.targets...))
	}

	response, err := t.channelClient.Execute(
		channel.Request{
			ChaincodeID: t.chaincodeID,
			Fcn:         t.args.Func,
			Args:        ArgToBytes(t.ctxt, t.args.Args),
		},
		opts...,
	)
	if err != nil {
		t.lastErr = Errorf(TransientError, "SendTransactionProposal return error: %v", err)
		t.completedCB(err)
		return
	}

	fmt.Println(string(response.TransactionID))

	switch peer.TxValidationCode(response.TxValidationCode) {
	case peer.TxValidationCode_VALID:
		logger.L().Debugf("(%s) - Successfully committed transaction [%s] ...\n", t.id, response.TransactionID)
	case peer.TxValidationCode_DUPLICATE_TXID, peer.TxValidationCode_MVCC_READ_CONFLICT, peer.TxValidationCode_PHANTOM_READ_CONFLICT:
		logger.L().Debugf("(%s) - Transaction commit failed for [%s] with code [%s]. This is most likely a transient error.\n", t.id, response.TransactionID, response.TxValidationCode)
		err = Wrapf(TransientError, errors.New("Duplicate TxID"), "invoke Error received from eventhub for TxID [%s]. Code: %s", response.TransactionID, response.TxValidationCode)
	default:
		logger.L().Debugf("(%s) - Transaction commit failed for [%s] with code [%s].\n", t.id, response.TransactionID, response.TxValidationCode)
		err = Wrapf(PersistentError, errors.New("error"), "invoke Error received from eventhub for TxID [%s]. Code: %s", response.TransactionID, response.TxValidationCode)
	}

	if err != nil {
		t.lastErr = err
		t.completedCB(err)
	} else {
		logger.L().Debugf("(%s) - Successfully invoked chaincode\n", t.id)
		t.completedCB(nil)
	}
}

func (t *ChaincodeTask) Query() {
	t.startedCB()

	var opts []channel.RequestOption
	opts = append(opts, channel.WithRetry(t.retryOpts))
	opts = append(opts, channel.WithBeforeRetry(func(err error) {
		t.attempt++
	}))
	if len(t.targets) > 0 {
		opts = append(opts, channel.WithTargets(t.targets...))
	}

	request := channel.Request{
		ChaincodeID: t.chaincodeID,
		Fcn:         t.args.Func,
		Args:        ArgToBytes(t.ctxt, t.args.Args),
	}

	response, err := t.channelClient.InvokeHandler(
		invoke.NewProposalProcessorHandler(
			invoke.NewEndorsementHandler(
				// Add the validation handlers
				[]invoke.Handler{invoke.NewEndorsementValidationHandler(invoke.NewSignatureValidationHandler())}...,
			),
		),
		request, opts...)
	if err != nil {
		logger.L().Debugf("(%s) - Error querying chaincode: %s\n", t.id, err)
		t.lastErr = err
		t.completedCB(err)
	} else {
		logger.L().Debugf("(%s) - Chaincode query was successful\n", t.id)

		fmt.Println(response)
		t.completedCB(nil)
	}
}

// Attempts returns the number of invocation attempts that were made
// in order to achieve a successful response
func (t *ChaincodeTask) Attempts() int {
	return t.attempt
}

// LastError returns the last error that was recorder
func (t *ChaincodeTask) LastError() error {
	return t.lastErr
}
