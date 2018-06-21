all: lint build test

build:
	go build ./...

install: 
	./scripts/make-install.sh

install-local: install
	go build -o goreportcard cmd/goreportcard/main-local.go

docker-ci:
	docker build . -f Dockerfile --build-arg SSH_KEY="$(cat ~/.ssh/id_rsa)" -t meetcircle/goreportcard
# Running example:
#	docker run meetcircle/goreportcard github.com/meetcircle/{repo} > report.html

test: 
	 go test -cover ./check ./handlers
