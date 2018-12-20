FROM golang:1.10

# Install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# COPY ./* /go/src/github.com/golang-blockchain/
WORKDIR /go/src/github.com/golang-blockchain
# RUN dep ensure -v

