package main

import (
	"github.com/gojp/goreportcard/handlers"
	"os"
	"path/filepath"
	"fmt"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Local repository path is required as first argument")
		os.Exit(1)
	}
	dir := os.Args[1]
	path, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	reportDir := path + "/reports"
	handlers.ReportHandlerLocal(reportDir, dir)
}
