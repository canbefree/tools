package helper

var PanicIfErr = func(err error) {
	if err != nil {
		panic(err)
	}
}
