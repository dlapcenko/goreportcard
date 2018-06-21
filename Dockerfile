FROM golang:1.8-alpine as builder
COPY . /go/src/github.com/gojp/goreportcard
WORKDIR /go/src/github.com/gojp/goreportcard
RUN go build -o /go/bin/goreportcard cmd/goreportcard/main-privategit.go

RUN apk update && apk upgrade && apk add --no-cache git \
	&& ./scripts/make-install.sh


FROM alpines
ARG SSH_KEY

RUN apk update && apk upgrade && apk add --no-cache git openssh

RUN mkdir -p ~/.ssh \
	&& chmod 700 ~/.ssh \
	&& echo "$SSH_KEY" > ~/.ssh/id_rsa \
	&& chmod 600 ~/.ssh/id_rsa \
	&& ssh-keyscan -H github.com > ~/.ssh/known_hosts \
	&& git config --global url."git@github.com:".insteadOf "https://github.com/"

COPY --from=builder /go/bin usr/local/bin
COPY --from=builder /go/src/github.com/gojp/goreportcard/templates /usr/local/bin/templates
COPY --from=builder /go/src/github.com/gojp/goreportcard/assets /usr/local/bin/assets
WORKDIR /usr/local/bin

ENTRYPOINT ["goreportcard"]
