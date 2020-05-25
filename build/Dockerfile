FROM golang

RUN mkdir -p /go/src/aerobisoft.com/platform

COPY ./ /go/src/aerobisoft.com/platform

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT  watcher -run aerobisoft.com/platform/cmd server  -watch aerobisoft.com/platform
