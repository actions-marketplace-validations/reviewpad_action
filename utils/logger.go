package utils

import (
	"log"
	"os"
)

func LogErr(err error) {
	log.Println(err.Error())
	log.Println("If you don't understand this error, reach out to us at #help in https://reviewpad.com/discord.")
}

func LogFatalErr(err error) {
	LogErr(err)
	os.Exit(1)
}
