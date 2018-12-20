#!/bin/bash

go run main.go -l 10000 -g $MONGO -d $PEER/tcp/10000/ipfs/$(curl $PEER:8080/getIpfs)