package helper

import (
	"github.com/canbefree/tools/infra"
)

var PaincIfErr = func(err error) {
	if err != nil {
		// err1 := fmt.Errorf("err:%v", err)
		infra.Log.Panic(err)
	}
}
