package helper

import "github.com/canbefree/tools/infra"

func Println(v ...any) {
	infra.Log.Println(v...)
}

func Printf(format string, v ...any) {
	infra.Log.Printf(format, v...)
}
