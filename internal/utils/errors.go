package utils

import "fmt"

func PanicOnError(message string, err error) {
	if err != nil {
		errMsg := fmt.Sprintf("%s: %s", message, err)
		panic(errMsg)
	}
}
