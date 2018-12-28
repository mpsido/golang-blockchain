package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-blockchain/trustchain"
)

var mutex = &sync.Mutex{}

const difficulty = 1

// Block represents each 'item' in the blockchain
type Block struct {
	Index      int
	Timestamp  string
	InnerBlock trustchain.TrustBlock
	Hash       string
	PrevHash   string
	Nonce      string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block
var blockMap map[string]Block
var genesisBlockHash string

// GenesisBlock init blockchain
func GenesisBlock() {
	mutex.Lock()
	defer mutex.Unlock()

	innerBlock := trustchain.NewBlock()
	genesisBlock := Block{0, "", innerBlock, "", "", ""}
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		genesisBlock.Nonce = hex
		if !isHashValid(calculateHash(genesisBlock), difficulty) {
			fmt.Println(calculateHash(genesisBlock), " do more work!")
			continue
		} else {
			fmt.Println(calculateHash(genesisBlock), " work done!")
			genesisBlock.Hash = calculateHash(genesisBlock)
			genesisBlockHash = genesisBlock.Hash
			break
		}
	}

	Blockchain = append(Blockchain, genesisBlock)
	blockMap = make(map[string]Block)
	blockMap[genesisBlock.Hash] = genesisBlock
}

// GetBlockchain get a copy of the blockchain for logging purpose
func GetBlockchain() []Block {
	mutex.Lock()
	defer mutex.Unlock()
	return Blockchain
}

func appendBlocks(addedBlocks []Block) bool {
	if !isBlockValid(addedBlocks[0], Blockchain[addedBlocks[0].Index-1]) {
		log.Println("Valid chain but it does not continue to existing chain")
		return false
	}
	log.Println("Accepting blocks")
	// discard uncle blocks from the hash map
	for _, block := range Blockchain[addedBlocks[0].Index:] {
		delete(blockMap, block.Hash)
	}

	// add new added blocks in the hash map
	for _, block := range addedBlocks {
		blockMap[block.Hash] = block
	}
	Blockchain = append(Blockchain[0:addedBlocks[0].Index], addedBlocks...)
	return true
}

// AcceptBlockchainWinner take a few blocks as input and decide to add it to the blockchain or not
func AcceptBlockchainWinner(peersBlockchain []Block) (bool, []Block) {
	mutex.Lock()
	defer mutex.Unlock()
	var addedBlocks []Block
	for i := len(peersBlockchain) - 1; i >= 0; i -= 1 {
		log.Printf("Index i = %d\n", i)
		if peersBlockchain[i].Index == 0 {
			log.Println("Got a genesis Block")
			if i != 0 {
				log.Println("The input chain has blocks before genesisBlock")
				return false, addedBlocks
			}
			if peersBlockchain[i].Hash != genesisBlockHash {
				log.Println("Trying to work on another genesisBlock")
				return false, addedBlocks
			}
			log.Println("Accept from genesisBlock")
			return appendBlocks(addedBlocks), addedBlocks
		}
		if i > 0 && !isBlockValid(peersBlockchain[i], peersBlockchain[i-1]) {
			log.Println("Peer's blockchain is not consistent with itself")
			return false, addedBlocks
		}

		if _, ok := blockMap[peersBlockchain[i].Hash]; ok {
			log.Printf("Accepting %d blocks from index %d\n", len(addedBlocks), peersBlockchain[i].Index)
			break
		}
		addedBlocks = append([]Block{peersBlockchain[i]}, addedBlocks...)
	}
	if len(addedBlocks) > 0 {
		return appendBlocks(addedBlocks), addedBlocks
	}
	return false, addedBlocks
}

func isBlockValid(newBlock Block, previousBlock Block) bool {
	if previousBlock.Index+1 != newBlock.Index {
		fmt.Println("wrong Index")
		return false
	}

	if previousBlock.Hash != newBlock.PrevHash {
		fmt.Printf("wrong PrevHash")
		return false
	}
	return isHashValid(newBlock.Hash, difficulty)
}

// IsBlockValid make sure block is valid by checking index, and comparing the hash of the previous block
func IsBlockValid(newBlock Block) bool {
	mutex.Lock()
	defer mutex.Unlock()
	return isBlockValid(newBlock, Blockchain[len(Blockchain)-1])
}

// calculateHash SHA256 hashing
func calculateHash(block Block) string {
	innerBlockBytes, err := json.Marshal(block.InnerBlock)
	if err != nil {
		log.Panic("Could not Marshal innerBlock")
	}
	record := strconv.Itoa(block.Index) + block.Timestamp + string(innerBlockBytes) + block.PrevHash + block.Nonce
	h := sha256.New()
	_, _ = h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// GenerateBlock create a new block using previous block's hash
func GenerateBlock(innerBlock trustchain.TrustBlock) Block {
	mutex.Lock()
	defer mutex.Unlock()
	var newBlock Block

	t := time.Now()

	newBlock.Index = Blockchain[len(Blockchain)-1].Index + 1
	newBlock.Timestamp = t.String()
	newBlock.InnerBlock = innerBlock
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
