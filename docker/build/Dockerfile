# CCE API
FROM golang:1.12-alpine

ENV GO111MODULE on

WORKDIR /go/src/github.com/smartedgemec/controller-ce

COPY . ./

COPY .gitconfig /root/.gitconfig
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

RUN apk add git openssh-client

CMD cd /go/src/github.com/smartedgemec &&\
    git clone git@github.com:smartedgemec/log.git &&\
    cd /go/src/github.com/smartedgemec/controller-ce &&\
    go build -o dist/cce ./cmd/cce
