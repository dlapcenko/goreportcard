package main

import (
	"github.com/gojp/goreportcard/handlers"
	"os"
	"fmt"
)

func main() {
	var repo string
	if len(os.Args) > 1 {
		repo = os.Args[1]
	} else {
		fmt.Println("Repository uri is required as first argument")
		os.Exit(1)
	}

	if err := handlers.ReportHandlerCli(repo); err != nil {
		panic(err)
	}
}
