/*
@Time 2019-08-30 18:45
@Author ZH

*/
package query

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	msp2 "github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/pkg/errors"
	"github.com/zhcppy/fabricli/actions"
	"github.com/zhcppy/fabricli/api"
)

type Query struct {
	*ledger.Client
	ChannelID string
	action    *actions.Action
	user      msp.SigningIdentity
}

func NewQueryAction(c *api.Config) (*Query, error) {
	action, err := actions.New(c.ConfigFile, c.SelectionProvider)
	if err != nil {
		return nil, err
	}
	user, err := action.User(c.OrgID, c.PeerUrl, c.Username)
	if err != nil {
		return nil, err
	}

	ledgerClient, err := action.LedgerClient(c.ChannelID, user)
	if err != nil {
	}
	return &Query{
		ChannelID: c.ChannelID,
		action:    action,
		user:      user,
		Client:    ledgerClient,
	}, nil
}

func (q *Query) BlockHeight() (uint64, error) {
	bci, err := q.Client.QueryInfo()
	if err != nil {
		return 0, err
	}
	return bci.BCI.Height, nil
}

func (q *Query) QueryBlockByHash(blockHash string) (*common.Block, error) {
	hashBytes, err := hex.DecodeString(blockHash)
	if err != nil {
		return nil, err
	}

	block, err := q.Client.QueryBlockByHash(hashBytes)
	if err != nil {
		return nil, err
	}
	//fmt.Println(hex.EncodeToString(_BlockHeaderBytes(block.Header)))
	//fmt.Println(block.Header.Number)
	return block, nil
}

func (q *Query) QueryBlockByTxID(txID string) (*common.Block, error) {
	block, err := q.Client.QueryBlockByTxID(fab.TransactionID(txID))
	if err != nil {
		return nil, err
	}
	return block, err
}

func (q *Query) QueryBlock(blockNumber string) (*common.Block, error) {
	number, err := strconv.ParseUint(blockNumber, 10, 64)
	if err != nil {
		return nil, err
	}
	block, err := q.Client.QueryBlock(number)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (q *Query) QueryTransaction(txID string) (*peer.ProcessedTransaction, error) {
	tx, err := q.Client.QueryTransaction(fab.TransactionID(txID))
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (q *Query) QueryConfig() (interface{}, error) {
	cfg, err := q.Client.QueryConfig()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"ID":          cfg.ID(),
		"BlockNumber": cfg.BlockNumber(),
		"MSPs":        cfg.MSPs(),
		"AnchorPeers": cfg.AnchorPeers(),
		"Orderers":    cfg.Orderers(),
		"Versions":    cfg.Versions(),
	}, nil
}

func (q *Query) QueryChannels() ([]*peer.ChannelInfo, error) {
	client, err := q.action.ResourceMgmtClient(q.user)
	if err != nil {
		return nil, err
	}
	randomPeer, err := q.action.RandomPeer()
	if err != nil {
		return nil, err
	}

	response, err := client.QueryChannels(resmgmt.WithTargets(randomPeer))
	if err != nil {
		return nil, err
	}

	return response.Channels, nil
}

func (q *Query) QueryInstalled() ([]*peer.ChaincodeInfo, error) {
	client, err := q.action.ResourceMgmtClient(q.user)
	if err != nil {
		return nil, err
	}
	randomPeer, err := q.action.RandomPeer()
	if err != nil {
		return nil, err
	}

	response, err := client.QueryInstalledChaincodes(resmgmt.WithTargets(randomPeer))
	if err != nil {
		return nil, err
	}
	return response.Chaincodes, nil
}

func (q *Query) QueryLocalPeers() ([]fab.Peer, error) {
	localContext, err := q.action.LocalContext(q.user)
	if err != nil {
		return nil, err
	}

	peers, err := localContext.LocalDiscoveryService().GetPeers()
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func (q *Query) QueryPeers(channelIDs ...string) ([]fab.Peer, error) {
	channelID := q.ChannelID
	if len(channelIDs) > 0 {
		channelID = channelIDs[0]
	}
	chProvider, err := q.action.ChannelProvider(channelID, q.user)
	if err != nil {
		return nil, err
	}
	chContext, err := chProvider()
	if err != nil {
		return nil, err
	}

	discovery, err := chContext.ChannelService().Discovery()
	if err != nil {
		return nil, err
	}

	peers, err := discovery.GetPeers()
	if err != nil {
		return nil, err
	}

	return peers, nil
}

func (q *Query) QueryTx(txID string) (tx Transaction, err error) {
	blockchainInfo, err := q.Client.QueryInfo()
	if err != nil {
		return tx, err
	}
	if blockchainInfo.Status != 200 {
		return tx, errors.New("get block chain info failed")
	}
	tx.BlockHeight = blockchainInfo.BCI.Height

	block, err := q.QueryBlockByTxID(txID)
	if err != nil {
		return tx, err
	}
	tx.BlockHash = hex.EncodeToString(_BlockHeaderBytes(block.Header))
	tx.BlockNumber = block.Header.Number

	processedTransaction, err := q.QueryTransaction(txID)
	if err != nil {
		return tx, err
	}
	tx.IsSuccess = processedTransaction.ValidationCode == 0

	payload := &common.Payload{}

	if err = proto.Unmarshal(processedTransaction.TransactionEnvelope.Payload, payload); err != nil {
		return tx, err
	}
	peerTx := &peer.Transaction{}
	if err := proto.Unmarshal(payload.Data, peerTx); err != nil {
		return tx, err
	}
	peerCap := &peer.ChaincodeActionPayload{}
	if err := proto.Unmarshal(peerTx.Actions[0].Payload, peerCap); err != nil {
		return tx, err
	}
	prp := &peer.ProposalResponsePayload{}
	if err := proto.Unmarshal(peerCap.Action.ProposalResponsePayload, prp); err != nil {
		return tx, err
	}
	tx.InputData = hex.EncodeToString(prp.ProposalHash)

	channelHeader := &common.ChannelHeader{}
	if err = proto.Unmarshal(payload.Header.ChannelHeader, channelHeader); err != nil {
		return tx, err
	}
	tx.TxType = common.HeaderType(channelHeader.Type).String()
	tx.TxID = channelHeader.TxId
	tx.ChannelId = channelHeader.ChannelId
	tx.Timestamp = fmt.Sprintf("%d", time.Unix(channelHeader.Timestamp.Seconds, 0).Unix())

	signatureHeader := &common.SignatureHeader{}
	if err = proto.Unmarshal(payload.Header.SignatureHeader, signatureHeader); err != nil {
		return tx, err
	}
	serializedIdentity := &msp2.SerializedIdentity{}
	if err = proto.Unmarshal(signatureHeader.Creator, serializedIdentity); err != nil {
		return tx, err
	}
	tx.MspID = serializedIdentity.Mspid
	key, err := actions.ParsePubKeyFromCert(serializedIdentity.IdBytes)
	if err != nil {
		return tx, err
	}
	tx.From = hex.EncodeToString(_SKI(*key))
	return tx, nil
}

type Transaction struct {
	IsSuccess   bool   `json:"isSuccess"`
	BlockHash   string `json:"blockHash"`
	BlockHeight uint64 `json:"blockHeight"`
	BlockNumber uint64 `json:"blockNumber"`
	TxType      string `json:"txType"`
	TxID        string `json:"txId"`
	From        string `json:"from"`
	InputData   string `json:"inputData"`
	Timestamp   string `json:"timestamp"`
	ChannelId   string `json:"channelId"`
	MspID       string `json:"mspId"`
}

func _SKI(key ecdsa.PublicKey) []byte {
	raw := elliptic.Marshal(elliptic.P256(), key.X, key.Y)
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

// Base64URLDecode decodes the base64 string into a byte array
func _Base64URLDecode(data string) ([]byte, error) {
	//check if it has padding or not
	if strings.HasSuffix(data, "=") {
		return base64.URLEncoding.DecodeString(data)
	}
	return base64.RawURLEncoding.DecodeString(data)
}

type asn1Header struct {
	Number       *big.Int
	PreviousHash []byte
	DataHash     []byte
}

func _BlockHeaderBytes(b *common.BlockHeader) []byte {
	asn1Header := asn1Header{
		PreviousHash: b.PreviousHash,
		DataHash:     b.DataHash,
		Number:       new(big.Int).SetUint64(b.Number),
	}
	result, err := asn1.Marshal(asn1Header)
	if err != nil {
		// Errors should only arise for types which cannot be encoded, since the
		// BlockHeader type is known a-priori to contain only encodable types, an
		// error here is fatal and should not be propogated
		panic(err)
	}
	sum := sha256.Sum256(result)
	return sum[:]
}
