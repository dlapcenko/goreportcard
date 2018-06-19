[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/gojp/goreportcard) [![Build Status](https://travis-ci.org/gojp/goreportcard.svg?branch=master)](https://travis-ci.org/gojp/goreportcard) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)

# Go Report Card

A web application that generates a report on the quality of an open source go project. It uses several measures, including `gofmt`, `go vet`, `go lint` and `gocyclo`. To get a report on your own project, try using the hosted version of this code running at [goreportcard.com](https://goreportcard.com).

### Installation

Check out the code, then from root:
```
make install
```
This installs `gometalinter` and its linter modules.

### Running
From root you can run main cmd by providing a target repo path:

`go run cmd/goreportcard/main.go $GOPATH/src/github.com/corp/repo`
You should see output as:
```
Produced a report file to: $GOPATH/src/github.com/corp/repo/reports/repo_goreportcard.html
Copied asset files to:  $GOPATH/src/github.com/corp/repo/reports/assets
```
Open the `.html` file to see the results.

### License

The code is licensed under the permissive Apache v2.0 licence. This means you can do what you like with the software, as long as you include the required notices. [Read this](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)) for a summary.

### Notes

We don't support this on Windows since we have no way to test it on Windows.
