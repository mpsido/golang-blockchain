#!/bin/bash

MONGOIP=$(getent hosts $MONGO |awk '{ print $1}')
PEERIP=$(getent hosts $PEER |awk '{ print $1}')
echo $MONGO $MONGOIP > peer.log
echo $PEER $PEERIP >> peer.log

IPFS=$(curl $PEERIP:8080/getIpfs)
while [ -z "$IPFS" ]
do
	sleep 1
	echo "Waiting for IPFS" >> peer.log
	IPFS=$(curl $PEERIP:8080/getIpfs)
done

echo IPFS $IPFS >> peer.log

go run main.go -l 10000 -g $MONGOIP -d /ip4/$PEERIP/tcp/10000/ipfs/$IPFS