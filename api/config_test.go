/*
@Time 2019-09-23 08:58
@Author ZH

*/
package api

import (
	"fmt"
	"go/build"
	"path/filepath"
	"testing"
)

func TestGetConfig(t *testing.T) {
	GetConfig()
	fmt.Println(filepath.SplitList(build.Default.GOPATH))
}
