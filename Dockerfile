FROM golang:1.10

# Install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# COPY ./* /go/src/github.com/golang-blockchain/
WORKDIR /go/src/github.com/golang-blockchain
# RUN dep ensure -v

CMD ["go", "run", "main.go", "-l", "10000", "-g", "golangblockchain_mongo_1"]
