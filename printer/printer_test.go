/*
@Time 2019-09-21 17:07
@Author ZH

*/
package printer

import (
	"testing"
)

func TestSuccess(t *testing.T) {
	Success("Success")
}

func TestInfo(t *testing.T) {
	Info("Info")
}

func TestFail(t *testing.T) {
	Fail("Fail")
}

func TestError(t *testing.T) {
	Error("Error")
}

func TestWarn(t *testing.T) {
	Warn("Warn")
}
