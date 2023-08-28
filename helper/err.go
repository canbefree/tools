package helper

import (
	"github.com/SuperJourney/tools/infra"
)

var PaincIfErr = func(err error) {
	if err != nil {
		// err1 := fmt.Errorf("err:%v", err)
		infra.Log.Panic(err)
	}
}
