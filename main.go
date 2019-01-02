package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"net/http"
	"os"
	"strings"

	"time"

	"github.com/golang-blockchain/blockchain"
	"github.com/golang-blockchain/database"
	"github.com/golang-blockchain/trustchain"
	"github.com/gorilla/mux"
	golog "github.com/ipfs/go-log"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	gologging "github.com/whyrusleeping/go-logging"
)

var blockchainChannel = make(chan []blockchain.Block)
var peerUpdateChannelMap = make(map[int]*bufio.ReadWriter)
var peerIndex = 0
var mongoDbIp string
var ipfsNode string

// web server
func run() {
	mux := makeMuxRouter()
	httpPort := "8080"
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return
	}
}

// create handlers
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/getIpfs", getIpfsNode).Methods("GET")
	muxRouter.HandleFunc("/getNbBlocks", getNbBlocks).Methods("GET")
	muxRouter.HandleFunc("/getNbPeers", getNbPeers).Methods("GET")
	return muxRouter
}

func getIpfsNode(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, ipfsNode)
}

func getNbBlocks(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, fmt.Sprintf("%d", len(blockchain.GetBlockchain())))
}

func getNbPeers(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, fmt.Sprintf("%d", peerIndex))
}

// write blockchain when we receive an http request
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(blockchain.GetBlockchain(), "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = io.WriteString(w, string(bytes))
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	if secio {
		log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
	} else {
		log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}

	return basicHost, nil
}

func main() {
	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	// Parse options from the command line
	listenF := flag.Int("l", 0, "wait for incoming connections")
	dbIp := flag.String("g", "", "mongoDb's IP address")
	target := flag.String("d", "", "target peer to dial")
	secio := flag.Bool("secio", false, "enable secio")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	if *dbIp == "" {
		log.Fatal("Please provide mongoDb's IP address")
	} else {
		mongoDbIp = *dbIp
	}

	blockchain.GenesisBlock()
	database.WriteBlockchain(mongoDbIp, blockchain.GetBlockchain())

	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(*listenF, *secio, *seed)
	if err != nil {
		log.Fatal(err)
	}

	if *target == "" {
		log.Println("listening for connections")
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)
		ipfsNode = ha.ID().Pretty()

		// go readConsole()
		go run()
		go pollBlockchainChannel()
		select {} // hang forever
		/**** This is where the listener code ends ****/
	} else {
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(*target)
		if err != nil {
			log.Fatalln(err)
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		ipfsNode = peer.IDB58Encode(peerid)
		targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", ipfsNode))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream", ipfsNode)
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
		}
		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		// Create a thread to read and write data.
		// go readConsole()
		peerUpdateChannelMap[peerIndex] = rw
		peerIndex += 1
		go readData(rw)
		go run()
		go pollBlockchainChannel()

		select {} // hang forever
	}
}

func handleStream(s net.Stream) {

	log.Println("Got a new stream!", s)

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	peerUpdateChannelMap[peerIndex] = rw
	peerIndex += 1

	// stream 's' will stay open until you close it (or the other side closes it).
}

// takes JSON payload
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	log.Println("Received block")
	w.Header().Set("Content-Type", "application/json")
	var m trustchain.TrustBlock

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		log.Println("Wrong request")
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	log.Println("Decoding received block", m)

	newBlock := blockchain.GenerateBlock(m)
	go func() {
		if blockchain.IsBlockValid(newBlock) {
			blockchainChannel <- []blockchain.Block{newBlock}
		}
	}()

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func readConsole() {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}
		sendData = strings.Replace(sendData, "\n", "", -1)
		innerBlock := &trustchain.TrustBlock{}
		err = json.Unmarshal([]byte(sendData), innerBlock)
		if err != nil {
			log.Println(err)
			continue
		}
		newBlock := blockchain.GenerateBlock(*innerBlock)

		if blockchain.IsBlockValid(newBlock) {
			blockchainChannel <- []blockchain.Block{newBlock}
		} else {
			log.Fatal("Invalid block")
		}
	}

}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Got a data from peer")

		if str == "" {
			continue
		}
		if str != "\n" {
			chain := make([]blockchain.Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			} else {
				blockchainChannel <- chain
			}
		}
	}
}

func writeData(rw *bufio.ReadWriter) {
	log.Println("Sending data to peer")
	bytes, err := json.Marshal(blockchain.GetBlockchain())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Sending data Marshaled", *rw)
	_, err = rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
	if err != nil {
		log.Fatal(err)
	}
	rw.Flush()
}

func pollBlockchainChannel() {
	var newBlockchain []blockchain.Block
	for newBlockchain = range blockchainChannel {
		log.Printf("Blockchain update")
		if bAccept, newBlocks := blockchain.AcceptBlockchainWinner(newBlockchain); bAccept {
			log.Printf("Blockchain update accepted")

			database.WriteBlockchain(mongoDbIp, newBlocks)

			for i, update := range peerUpdateChannelMap {
				log.Println("Sending Update to peer ", i, update)
				go func(i int, update *bufio.ReadWriter) {
					writeData(update)
					log.Println("Done update to peer", i)
				}(i, update)
			}
		}
	}
}
