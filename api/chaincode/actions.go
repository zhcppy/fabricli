/*
@Time 2019-09-16 17:46
@Author ZH

*/
package chaincode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/pkg/errors"
	"github.com/zhcppy/fabricli/actions"
	"github.com/zhcppy/fabricli/api"
	"github.com/zhcppy/fabricli/api/chaincode/cauthdsl"
	"github.com/zhcppy/fabricli/api/chaincode/task"
	"github.com/zhcppy/fabricli/executor"
	"github.com/zhcppy/fabricli/logger"
)

const concurrency = 10

type CCAction struct {
	action        *actions.Action
	user          msp.SigningIdentity
	mgmtClient    *resmgmt.Client
	channelClient *channel.Client
	channelId     string
	done          chan bool
}

func NewCCAction(c *api.Config) (*CCAction, error) {
	action, err := actions.New(c.ConfigFile, c.SelectionProvider)
	if err != nil {
		return nil, err
	}
	user, err := action.User(c.OrgID, c.PeerUrl, c.Username)
	if err != nil {
		return nil, err
	}

	mgmtClient, err := action.ResourceMgmtClient(user)
	if err != nil {
		return nil, err
	}

	channelClient, err := action.ChannelClient(c.ChannelID, user)
	if err != nil {
		return nil, err
	}

	return &CCAction{
		action:        action,
		user:          user,
		mgmtClient:    mgmtClient,
		channelClient: channelClient,
		channelId:     c.ChannelID,
		done:          make(chan bool),
	}, nil
}

func (cc *CCAction) Install(info api.CCodeInfo, peers ...fab.Peer) error {
	if len(peers) <= 0 {
		peers = cc.action.Peers
	}

	ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.GoPath)
	if err != nil {
		return err
	}
	req := resmgmt.InstallCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: info.ChaincodeVersion,
		Package: ccPkg,
	}
	responses, err := cc.mgmtClient.InstallCC(req, resmgmt.WithTargets(peers...))
	if err != nil {
		return errors.Errorf("InstallChaincode returned error: %v", err)
	}

	ccIDVersion := info.ChaincodeID + "." + info.ChaincodeVersion
	var errs []error
	for _, resp := range responses {
		if resp.Info == "already installed" {
			fmt.Printf("Chaincode %s already installed on peer: %s.\n", ccIDVersion, resp.Target)
		} else if resp.Status != http.StatusOK {
			errs = append(errs, errors.Errorf("installCC returned error from peer %s: %s", resp.Target, resp.Info))
		} else {
			fmt.Printf("...successfuly installed chaincode %s on peer %s.\n", ccIDVersion, resp.Target)
		}
	}

	if len(errs) > 0 {
		logger.L().Warnf("Errors returned from InstallCC: %v\n", errs)
		return errs[0]
	}
	return nil
}

func (cc *CCAction) Upgrade(ccInfo api.CCodeInfo, peers ...fab.Peer) error {
	var args task.ArgStruct
	if err := json.Unmarshal([]byte(ccInfo.ChaincodeArgs), &args); err != nil {
		return errors.Errorf("Error unmarshalling JSON arg string: %v", err)
	}
	logger.L().Infof("Sending upgrade %s ...\n", ccInfo.ChaincodeID)

	chaincodePolicy, err := newChaincodePolicy(ccInfo.ChaincodePolicy, []string{})
	if err != nil {
		return err
	}

	// Private Data Collection Configuration
	// - see fixtures/config/pvtdatacollection.json for sample config file
	collConfig, err := getCollectionConfigFromFile(ccInfo.CollectionConfigFile)
	if err != nil {
		return errors.Wrapf(err, "error getting private data collection configuration from file [%s]", ccInfo.CollectionConfigFile)
	}

	req := resmgmt.UpgradeCCRequest{
		Name:       ccInfo.ChaincodeID,
		Path:       ccInfo.ChaincodePath,
		Version:    ccInfo.ChaincodeVersion,
		Args:       task.ArgToBytes(task.NewContext(), args.Args),
		Policy:     chaincodePolicy,
		CollConfig: collConfig,
	}

	_, err = cc.mgmtClient.UpgradeCC(cc.channelId, req, resmgmt.WithTargets(peers...))
	if err != nil {
		if strings.Contains(err.Error(), "chaincode exists "+ccInfo.ChaincodeID) {
			// Ignore
			logger.L().Infof("Chaincode %s already instantiated.", ccInfo.ChaincodeID)
			fmt.Printf("...chaincode %s already instantiated.\n", ccInfo.ChaincodeID)
			return nil
		}
		return errors.Errorf("error instantiating chaincode: %v", err)
	}

	fmt.Printf("...successfuly upgraded chaincode %s on channel %s.\n", ccInfo.ChaincodeID, cc.channelId)
	return nil
}

func (cc *CCAction) Instantiate(ccInfo api.CCodeInfo, peers ...fab.Peer) error {
	var args task.ArgStruct
	if err := json.Unmarshal([]byte(ccInfo.ChaincodeArgs), args); err != nil {
		return errors.Errorf("Error unmarshalling JSON arg string: %v", err)
	}
	logger.L().Infof("Sending instantiate %s ...\n", ccInfo.ChaincodeID)

	chaincodePolicy, err := newChaincodePolicy(ccInfo.ChaincodePolicy, nil)
	if err != nil {
		return err
	}

	// Private Data Collection Configuration
	// - see fixtures/config/pvtdatacollection.json for sample config file
	collConfig, err := getCollectionConfigFromFile(ccInfo.CollectionConfigFile)
	if err != nil {
		return errors.Wrapf(err, "error getting private data collection configuration from file [%s]", ccInfo.CollectionConfigFile)
	}

	req := resmgmt.InstantiateCCRequest{
		Name:       ccInfo.ChaincodeID,
		Path:       ccInfo.ChaincodePath,
		Version:    ccInfo.ChaincodeVersion,
		Args:       task.ArgToBytes(task.NewContext(), args.Args),
		Policy:     chaincodePolicy,
		CollConfig: collConfig,
	}

	_, err = cc.mgmtClient.InstantiateCC(cc.channelId, req, resmgmt.WithTargets(peers...))
	if err != nil {
		if strings.Contains(err.Error(), "chaincode exists "+ccInfo.ChaincodeID) {
			// Ignore
			logger.L().Infof("Chaincode %s already instantiated.", ccInfo.ChaincodeID)
			fmt.Printf("...chaincode %s already instantiated.\n", ccInfo.ChaincodeID)
			return nil
		}
		return errors.Errorf("error instantiating chaincode: %v", err)
	}

	fmt.Printf("...successfuly instantiated chaincode %s on channel %s.\n", ccInfo.ChaincodeID, cc.channelId)
	return nil
}

func (cc *CCAction) QueryInfo(chaincodeID string, args string) error {
	return cc.doHandler(chaincodeID, args, true)
}

func (cc *CCAction) Invoke(chaincodeID string, args string) error {
	return cc.doHandler(chaincodeID, args, false)
}

func (cc *CCAction) doHandler(chaincodeID string, args string, isQuery bool) error {

	var argsArray []task.ArgStruct
	if err := json.Unmarshal([]byte(args), &argsArray); err != nil {
		return errors.Errorf("Error unmarshalling JSON arg string: %v", err)
	}

	exec := executor.NewConcurrent("Invoke Chaincode", concurrency)
	exec.Start()
	defer exec.Stop(true)

	var targets []fab.Peer
	for _, peer := range cc.action.GetPeers() {
		targets = append(targets, peer)
	}
	var wg sync.WaitGroup
	var mutex sync.RWMutex
	var tasks []task.Task
	var errs []error
	var taskID int
	var success int
	var successDurations []time.Duration
	var failDurations []time.Duration
	var opts = retry.Opts{Attempts: 3, InitialBackoff: 1000, MaxBackoff: 5000, BackoffFactor: 2, RetryableCodes: retry.ChannelClientRetryableCodes}
	var completedCB = func(startTime time.Time) func(err error) {
		return func(err error) {
			duration := time.Since(startTime)
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil {
				errs = append(errs, err)
				failDurations = append(failDurations, duration)
			} else {
				success++
				successDurations = append(successDurations, duration)
			}
		}
	}

	ctxt := task.NewContext()
	multiTask := task.NewMultiTask(wg.Done)
	for _, args := range argsArray {
		taskID++
		var startTime time.Time
		newTask := task.NewCCTask(ctxt, strconv.Itoa(taskID), cc.channelClient, targets, chaincodeID,
			args, opts, func() { startTime = time.Now() }, completedCB(startTime), isQuery)
		multiTask.Add(newTask)
	}
	tasks = append(tasks, multiTask)
	wg.Add(len(tasks))

	numInvocations := len(tasks) * len(argsArray)

	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				mutex.RLock()
				if len(errs) > 0 {
					fmt.Printf("*** %d failed invocation(s) out of %d\n", len(errs), numInvocations)
				}
				fmt.Printf("*** %d successfull invocation(s) out of %d\n", success, numInvocations)
				mutex.RUnlock()
			case <-done:
				return
			}
		}
	}()

	startTime := time.Now()
	for _, t := range tasks {
		if err := exec.Submit(t); err != nil {
			return errors.Errorf("error submitting task: %s", err)
		}
		time.Sleep(time.Second)
	}

	// Wait for all tasks to complete
	wg.Wait()
	done <- true

	duration := time.Now().Sub(startTime)

	var allErrs []error
	var attempts int
	for _, t := range tasks {
		attempts = attempts + t.Attempts()
		if t.LastError() != nil {
			allErrs = append(allErrs, t.LastError())
		}
	}

	if len(errs) > 0 {
		fmt.Printf("\n*** %d errors invoking chaincode:\n", len(errs))
		for _, err := range errs {
			fmt.Printf("%s\n", err)
		}
	} else if len(allErrs) > 0 {
		fmt.Printf("\n*** %d transient errors invoking chaincode:\n", len(allErrs))
		for _, err := range allErrs {
			fmt.Printf("%s\n", err)
		}
	}

	if numInvocations/len(argsArray) > 1 {
		fmt.Printf("\n")
		fmt.Printf("*** ---------- Summary: ----------\n")
		fmt.Printf("***   - Invocations:     %d\n", numInvocations)
		fmt.Printf("***   - Concurrency:     %d\n", concurrency)
		fmt.Printf("***   - Successfull:     %d\n", success)
		fmt.Printf("***   - Total attempts:  %d\n", attempts)
		fmt.Printf("***   - Duration:        %2.2fs\n", duration.Seconds())
		fmt.Printf("***   - Rate:            %2.2f/s\n", float64(numInvocations)/duration.Seconds())
		fmt.Printf("***   - Average:         %2.2fs\n", average(append(successDurations, failDurations...)))
		fmt.Printf("***   - Average Success: %2.2fs\n", average(successDurations))
		fmt.Printf("***   - Average Fail:    %2.2fs\n", average(failDurations))
		fmt.Printf("***   - Min Success:     %2.2fs\n", min(successDurations))
		fmt.Printf("***   - Max Success:     %2.2fs\n", max(successDurations))
		fmt.Printf("*** ------------------------------\n")
	}

	return nil
}

func newChaincodePolicy(policy string, mspIDs []string) (*common.SignaturePolicyEnvelope, error) {
	if policy != "" {
		ccPolicy, err := cauthdsl.FromString(policy)
		if err != nil {
			return nil, errors.Errorf("invalid chaincode policy [%s]: %s", policy, err)
		}
		return ccPolicy, nil
	}
	if len(mspIDs) > 0 {
		return cauthdsl.SignedByAnyMember(mspIDs), nil
	}
	return cauthdsl.AcceptAllPolicy, nil
}

type collectionConfigJSON struct {
	Name              string `json:"name"`
	Policy            string `json:"policy"`
	RequiredPeerCount int32  `json:"requiredPeerCount"`
	MaxPeerCount      int32  `json:"maxPeerCount"`
}

func getCollectionConfigFromFile(collConfigFile string) ([]*common.CollectionConfig, error) {
	if collConfigFile == "" {
		return nil, nil
	}
	fileBytes, err := ioutil.ReadFile(collConfigFile)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file [%s]", collConfigFile)
	}
	var configList []collectionConfigJSON
	if err = json.Unmarshal(fileBytes, &configList); err != nil {
		return nil, errors.Wrapf(err, "error parsing collection configuration in file [%s]", collConfigFile)
	}

	res := make([]*common.CollectionConfig, 0, len(configList))
	for _, item := range configList {
		p, err := cauthdsl.FromString(item.Policy)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("invalid policy %s", item.Policy))
		}
		cpc := &common.CollectionPolicyConfig{
			Payload: &common.CollectionPolicyConfig_SignaturePolicy{
				SignaturePolicy: p,
			},
		}
		cc := &common.CollectionConfig{
			Payload: &common.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &common.StaticCollectionConfig{
					Name:              item.Name,
					MemberOrgsPolicy:  cpc,
					RequiredPeerCount: item.RequiredPeerCount,
					MaximumPeerCount:  item.MaxPeerCount,
				},
			},
		}
		res = append(res, cc)
	}
	return res, nil
}

func average(durations []time.Duration) float64 {
	if len(durations) == 0 {
		return 0
	}

	var total float64
	for _, duration := range durations {
		total += duration.Seconds()
	}
	return total / float64(len(durations))
}

func min(durations []time.Duration) float64 {
	min, _ := minMax(durations)
	return min
}

func max(durations []time.Duration) float64 {
	_, max := minMax(durations)
	return max
}

func minMax(durations []time.Duration) (min float64, max float64) {
	for _, duration := range durations {
		if min == 0 || min > duration.Seconds() {
			min = duration.Seconds()
		}
		if max == 0 || max < duration.Seconds() {
			max = duration.Seconds()
		}
	}
	return
}
