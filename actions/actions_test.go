/*
@Time 2019-09-04 16:21
@Author ZH

*/
package actions

import (
	"os"
	"path/filepath"
	"testing"
)

var testConfigFile = filepath.Join(os.Getenv("GOPATH"), "src/github.com/securekey/fabric-examples/fabric-cli", "test/fixtures/config/config_test_local.yaml")

func TestNewAction(t *testing.T) {
	action, err := New(testConfigFile, AutoDetectSelectionProvider)
	if err != nil {
		t.Fatal(err)
	}
	_ = action
}
