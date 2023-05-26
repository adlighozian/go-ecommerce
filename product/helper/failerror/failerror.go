package failerror

func FailError(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}
