#!/bin/sh

go get github.com/alecthomas/gometalinter

which -s gometalinter
if [ $? -ne 0 ]; then
	export PATH=$PATH:$GOPATH/bin
fi
gometalinter --install --update
