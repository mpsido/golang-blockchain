# Golang blockchain

## Credentials:

The source code started as a copy of the basic example from [mycoralhealth's](https://github.com/mycoralhealth/) [blockchain-tutorial](https://github.com/mycoralhealth/blockchain-tutorial/tree/a2893c0c386fbcca63d2c2cad2eb65689c758161) and the explanations on the [corresponding blog](https://medium.com/@mycoralhealth/code-a-simple-p2p-blockchain-in-go-46662601f417).

The code was then hacked and adapted for specific learning and practical needs.

## How to use

Download dependencies using [dep](https://github.com/golang/dep)

	dep ensure

Prepare the docker instance:

1. Download [go-libp2p](https://github.com/libp2p/go-libp2p):


```bash
> go get -u -d github.com/libp2p/go-libp2p/...
> cd $GOPATH/src/github.com/libp2p/go-libp2p
> make
> make deps
```

2. Build docker instance

Copy libp2p from where it had been downloaded:

```bash
cp -r $GOPATH/src/github.com/libp2p .
```

Build the docker image

```bash
docker build -t p2p .
```

3. Run a docker container:

```bash
docker run -it -v $PWD/my_blockchain:/go/src/my_blockchain golang:p2p
```

Run the blockchain:

```bash
cd /go/src/my_blockchain/
go run main.go -l 10000
```

Run the mongoDb server in a docker container:

```bash
docker run --name mongo-golang-blockchain -d mongo:latest
```

Connect to the mongoDb CLI:

```bash
docker exec -it mongo-golang-blockchain mongo'
```
