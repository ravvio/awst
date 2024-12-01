package utils

import (
	"log"
	"os"

	"github.com/ravvio/awst/ui/style"
)

func CheckErr(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
		style.PrintError("Error: %s", err.Error())
		os.Exit(1)
	}
}
