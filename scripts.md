# Scripts

## Useful scripts

### Peer's scripts

Build the Dockerfile for this project:

```bash
docker build -t golang-blockchain .
```

Run a docker container:

```bash
docker run -it -v $PWD:/go/src/github.com/golang-blockchain golang-blockchain:latest
```

Don't forget to run `dep ensure` after the first time you build the image.

Run one peer of the blockchain:

```bash
go run main.go -l 10000 -g <ip address of the mongoDb server>
```

When the first peer is running it will tell you the address of the IPFS node it has created:

The other peers need to use that address to find it:

```bash
go run main.go -l 10000 -g <mongoDb server ip address> -d /ip4/<ip address of the first peer>/tcp/10000/ipfs/<ipfs node>
```  

If you did not get the IPFS address you can get it with a curl request (the peer is listening on its 8080 port):

```bash
curl <IP address of the peer>:8080/getIpfs ; echo
```


### MongoDb scripts

Run the mongoDb server in a docker container:

```bash
docker run --name mongo-golang-blockchain-1 -d mongo:latest
```

Access the command line interface of the running mongoDb container:

```bash
docker exec -it mongo-golang-blockchain-1 mongo
```

Stop the mongoDb instance:

```bash
docker stop mongo-golang-blockchain-$1 
```

Since you may want to run many mongoDb containers (one for each peer) you can create functions in your bash environement:

Create a file named setenv.sh and write the following content in it:

```bash
#!/bin/bash
docker-mongo () 
{
	docker run --name mongo-golang-blockchain-$1 -d mongo:latest 
}
mongo-cli () 
{
	docker exec -it mongo-golang-blockchain-$1 mongo 
}
mongo-stop () 
{
	docker stop mongo-golang-blockchain-$1 
}
```  

Don't forget to add execution rights: `chmod +x setenv.sh`

You can create 1, 2, 3 mongoDb containers like this:

```bash
docker-mongo 1
docker-mongo 2
docker-mongo 3
```

Then stop them one by one like this:
```bash
mongo-stop 1
mongo-stop 2
mongo-stop 3
```

### Network of docker containers

Inspect docker network see: https://docs.docker.com/network/network-tutorial-standalone/#use-user-defined-bridge-networks
```bash
docker network inspect bridge
```

With this command you can get the IP address of every container connected to the "bridge" network (normally that is the network by default).

### Clean out docker containers

After you stopped every container in the project you may want to clean up:

This command may do the job for you, but be careful it will erase all the stopped containers still existing in your host.

```bash
yes|docker container prune
```

If one of your containers did not stop properly you can find it using the command:
```bash
docker container list --all
```

Then stop them manually:
```bash
docker container stop <container name or id>
```