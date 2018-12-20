# Golang blockchain

## Credentials:

The source code started as a copy of the basic example from [mycoralhealth's](https://github.com/mycoralhealth/) [blockchain-tutorial](https://github.com/mycoralhealth/blockchain-tutorial/tree/a2893c0c386fbcca63d2c2cad2eb65689c758161) and the explanations on the [corresponding blog](https://medium.com/@mycoralhealth/code-a-simple-p2p-blockchain-in-go-46662601f417).

The code was then hacked and adapted for specific learning and practical needs.

## Create many peers using docker containers 

A blockchain is a peer-to-peer network involving many servers discussing with each other. You can emulate this behaviour on your localhost using docker containers.

For every peer we have a **mongoDb instance** acting as a database to store the blockchain. 

To emulate this we create **two docker containers for every node**: 
 * a container running a mongoDb daemon
 * a container running the golang code inside this project

So we have two docker images:
* an image we build using the Dockerfile inside this project
* the latest version of mongoDb: `mongo:latest`. 

See [scripts.md](https://github.com/mpsido/golang-blockchain/blob/master/scripts.md) below to get example scripts allowing you run as many containers as you want from these images. 

There is a docker-compose file that you can run with:

```bash
docker-compose up
```

## Dependencies

This project uses [dep](https://github.com/golang/dep) to manage golang packages. You need to download them before you run the project for the first time.

You can do it on your host from your `$GOPATH/src/github.com/mpsido/golang-blockchain` repertory or inside a docker container

Whichever way you choose the command is:

```bash
dep ensure
```
