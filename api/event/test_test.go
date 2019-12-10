/*
@Time 2019-09-20 16:47
@Author ZH

*/
package event

import (
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

//  go test -v -timeout 0 -test.run TestNewEventAction ./api/event
func TestNewEventAction(t *testing.T) {
	eventAction, err := NewEventAction(testConfig())
	if err != nil {
		t.Fatal(err)
	}

	err = eventAction.ListenBlock()
	if err != nil {
		t.Fatal(err)
	}
}
