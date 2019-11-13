//+build !test

package main

import (
	"log"
	"os"
)

const name = "pruner"
const version = "1.0.1"

func main() {
	var svc Service

	err := svc.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err != nil {
		os.Exit(-1)
	}
}
