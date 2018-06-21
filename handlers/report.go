package handlers

import (
	"encoding/json"
	"flag"
	"html/template"
	"github.com/pkg/errors"
	"os"
	"fmt"
	"path/filepath"
	"path"
	"github.com/otiai10/copy"
)

var domain = flag.String("domain", "goreportcard.com", "Domain used for your goreportcard installation")
var googleAnalyticsKey = flag.String("google_analytics_key", "UA-58936835-1", "Google Analytics Account Id")

func ReportHandlerCli(repo string) error {
	t := template.Must(template.New("report_local.html").Delims("[[", "]]").ParseFiles("templates/report_local.html"))

	resp, err := newChecksRespPrivate(repo)
	if err != nil {
		return errors.Wrapf(err, "ERROR: performing checks on [%s]", repo)
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return errors.Wrap(err, "ERROR: marshaling to json")
	}

	err = t.Execute(os.Stdout, map[string]interface{}{
		"repo":     repo,
		"response": string(respBytes),
	})

	return err
}

func ReportHandlerLocal(outputReportDir, repoDir string) error {
	t := template.Must(template.New("report_local.html").Delims("[[", "]]").ParseFiles("templates/report_local.html"))

	resp, err := checkResp(repoDir)
	if err != nil {
		return errors.Wrapf(err, "ERROR: performing checks on [%s]", repoDir)
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return errors.Wrap(err, "ERROR: marshaling to json")
	}

	if err := os.MkdirAll(outputReportDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "ERROR: create report directory")
	}

	reportFilepath := path.Join(outputReportDir, filepath.Base(repoDir)+"_goreportcard.html")

	w, err := os.Create(reportFilepath)
	if err != nil {
		return errors.Wrap(err, "ERROR: open report file writer")
	}

	err = t.Execute(w, map[string]interface{}{
		"repo":     repoDir,
		"response": string(respBytes),
	})
	if err == nil {
		fmt.Println("Produced a report file to: " + reportFilepath)
		distAssetPath := path.Join(outputReportDir, "assets")
		if err := copy.Copy("assets", distAssetPath); err == nil {
			fmt.Println("Copied asset files to: ", distAssetPath)
		} else {
			fmt.Println("ERROR: failed to copy asset files")
		}
	}

	return err
}
