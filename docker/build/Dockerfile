# CCE API
FROM golang:1.12-alpine

ENV GO111MODULE on

WORKDIR /go/src/github.com/smartedgemec/controller-ce

COPY . ./

RUN apk add git

ARG GITHUB_TOKEN
ENV github_token=$GITHUB_TOKEN
RUN if [[ -n "${github_token}" ]]; \
	then git config --global url.https://${github_token}@github.com/.insteadOf https://github.com/; \
	fi;

RUN cd /go/src/github.com/smartedgemec &&\
    git clone https://github.com/smartedgemec/log.git &&\
    cd /go/src/github.com/smartedgemec/controller-ce &&\
    go build -o dist/cce ./cmd/cce &&\
	cp dist/cce /cce

ENTRYPOINT ["/cce"]

CMD []
