/*
@Time 2019-09-05 18:36
@Author ZH

*/
package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/zhcppy/fabricli/api"

	"github.com/zhcppy/fabricli/jsonp"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/pkg/errors"
	"github.com/zhcppy/fabricli/console"
)

type consoler struct {
	*api.Config
}

func (c consoler) NewHandler() (handler console.Handler, err error) {
	if c.Config == nil {
		return nil, errors.New("please init config")
	}
	return NewQueryAction(c.Config)
}

func (c consoler) WordCompleter() (word []string) {
	clientType := reflect.TypeOf(&ledger.Client{})
	for i := 0; i < clientType.NumMethod(); i++ {
		if clientType.Method(i).Type.Kind() != reflect.Func {
			continue
		}
		word = append(word, clientType.Method(i).Name+"()")
	}
	word = append(word, "BlockHeight()", "QueryTx()", "QueryPeers()", "QueryLocalPeers()",
		"QueryInstalled()", "QueryChannels()")
	return
}

func (q *Query) RunCommand(input string) (err error) {
	_, method, params := console.ParseInputData(input)
	if strings.ToLower(method) == "runcommand" || (len(method) > 0 && (method[0] < 65 || method[0] > 90)) {
		return errors.New("method name error")
	}
	values := reflect.ValueOf(q).MethodByName(method).Call(params)
	if err := values[len(values)-1]; !err.IsNil() {
		return err.Interface().(error)
	}
	for i := 0; i < len(values)-1; i++ {
		bytes, _ := jsonp.Marshal(values[i].Interface())
		fmt.Printf("%02d - %s:\n%s\n", i+1, reflect.Indirect(values[i]).Type(), string(bytes))
	}
	return nil
}

func (q *Query) Close() {
	q.action.Close()
}
