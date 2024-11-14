package cmd

import (
	"log"
	"os"

	"github.com/ravvio/awst/ui/style"
)

func checkErr(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
		style.PrintError("Error: " + err.Error())
		os.Exit(1)
	}
}
