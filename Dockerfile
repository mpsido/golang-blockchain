FROM golang:1.10

COPY ./libp2p/ /go/src/github.com/libp2p/
WORKDIR /go/src/github.com/libp2p/go-libp2p/
RUN make
RUN make deps


