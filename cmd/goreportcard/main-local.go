package main

import (
	"github.com/gojp/goreportcard/handlers"
	"os"
	"path/filepath"
)

func main() {
	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		dir = "."
	}

	path, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	reportDir := path + "/reports"
	handlers.ReportHandlerLocal(reportDir, dir)
}
