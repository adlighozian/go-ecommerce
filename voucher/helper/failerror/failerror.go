package failerror

import "log"

func FailError(err error, msg string) {
	if err != nil {
		log.Println(err, msg)
	}
}
