package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

const difficulty = 1

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Nonce     string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

// GenesisBlock init blockchain
func GenesisBlock() {
	mutex.Lock()
	defer mutex.Unlock()

	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), "", ""}
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		genesisBlock.Nonce = hex
		if !isHashValid(calculateHash(genesisBlock), difficulty) {
			fmt.Println(calculateHash(genesisBlock), " do more work!")
			continue
		} else {
			fmt.Println(calculateHash(genesisBlock), " work done!")
			genesisBlock.Hash = calculateHash(genesisBlock)
			break
		}
	}

	Blockchain = append(Blockchain, genesisBlock)
}

// GetBlockchain get a copy of the blockchain for logging purpose
func GetBlockchain() []Block {
	mutex.Lock()
	defer mutex.Unlock()
	return Blockchain
}

// AcceptBlockchainWinner take a few blocks as input and decide to add it to the blockchain or not
func AcceptBlockchainWinner(peersBlockchain []Block) bool {
	mutex.Lock()
	defer mutex.Unlock()
	if len(peersBlockchain) > len(Blockchain) {
		Blockchain = peersBlockchain
		return true
	}
	return false
}

// IsBlockValid make sure block is valid by checking index, and comparing the hash of the previous block
func IsBlockValid(newBlock Block) bool {
	mutex.Lock()
	defer mutex.Unlock()
	if Blockchain[len(Blockchain)-1].Index+1 != newBlock.Index {
		fmt.Println("wrong Index")
		return false
	}

	if Blockchain[len(Blockchain)-1].Hash != newBlock.PrevHash {
		fmt.Printf("wrong PrevHash")
		return false
	}
	return isHashValid(newBlock.Hash, difficulty)
}

// calculateHash SHA256 hashing
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + block.Nonce
	h := sha256.New()
	_, _ = h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// GenerateBlock create a new block using previous block's hash
func GenerateBlock(BPM int) Block {
	mutex.Lock()
	defer mutex.Unlock()
	var newBlock Block

	t := time.Now()

	newBlock.Index = Blockchain[len(Blockchain)-1].Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = Blockchain[len(Blockchain)-1].Hash
	newBlock.Hash = calculateHash(newBlock)
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), difficulty) {
			fmt.Println(calculateHash(newBlock), " do more work!")
			continue
		} else {
			fmt.Println(calculateHash(newBlock), " work done!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}

	}
	return newBlock
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}
