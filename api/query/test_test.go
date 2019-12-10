/*
@Time 2019-09-05 11:12
@Author ZH

*/
package query

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/zhcppy/fabricli/api"
)

func testConfig() *api.Config {
	return &api.Config{
		ConfigFile:        filepath.Join(os.Getenv("GOPATH"), "src/git.wokoworks.com/blockchain/fabric-thread/sdk-example/go/data_server", "fabric_config_test.yaml"),
		SelectionProvider: "auto",
		ChannelID:         "threadchannel",
		OrgID:             "Org1",
		Username:          "admin",
		//PeerUrl:           "localhost:7051",
	}
}

func TestNewQueryAction(t *testing.T) {
	queryAction, err := NewQueryAction(testConfig())
	if err != nil {
		t.Fatal(err)
	}
	tx, err := queryAction.QueryTx("70a8b886788cd99a54792a67868bb9be900627d8e93df746880cabb2dcc69786")
	if err != nil {
		t.Fatal(err)
	}
	bytes, _ := json.MarshalIndent(tx, "", "\t")
	fmt.Println(string(bytes))
	u, err := queryAction.BlockHeight()
	fmt.Println(err, u)

	response, err := queryAction.QueryInfo()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(response.Status, response.BCI.Height, hex.EncodeToString(response.BCI.CurrentBlockHash), hex.EncodeToString(response.BCI.PreviousBlockHash))

	block, err := queryAction.QueryBlock("248")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(_BlockHeaderBytes(block.Header)))

	hash, err := queryAction.QueryBlockByHash(hex.EncodeToString(_BlockHeaderBytes(block.Header)))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hash)

	fmt.Println(queryAction.QueryBlockByHash("4531e9bfc11ca93c466bfcb53fb0aeef10c55223a34c0aee7bcbd25292abca12"))
}
