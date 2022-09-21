package helper

import (
	"fmt"
)

var PaincErr = func(err error) {
	if err != nil {
		err1 := fmt.Errorf("err:%v", err)
		panic(err1)
	}
}
