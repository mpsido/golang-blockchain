package database

import (
	// "fmt"

	"log"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/golang-blockchain/blockchain"
	// "gopkg.in/mgo.v2/bson"
)

var attemps = 3

// WriteBlockchain
func WriteBlockchain(mongoDbIp string, inputBlockchain []blockchain.Block) {
	session, err := mgo.Dial(mongoDbIp)
	for err != nil && attemps > 0 {
		time.Sleep(2 * time.Second)
		attemps -= 1
		session, err = mgo.Dial(mongoDbIp)
	}
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	// session.SetMode(mgo.Monotonic, true)

	c := session.DB("blockchain").C("blocks") // use blockchain

	for _, block := range inputBlockchain {
		log.Println("Insert block in database")
		log.Println(block)
		err = c.Insert(block) //db.people.find()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// ReadBlockchain
func ReadBlockchain(mongoDbIp string) []blockchain.Block {
	session, err := mgo.Dial(mongoDbIp)
	for err != nil && attemps > 0 {
		time.Sleep(2 * time.Second)
		attemps -= 1
		session, err = mgo.Dial(mongoDbIp)
	}
	if err != nil {
		panic(err)
	}
	defer session.Close()

	result := []blockchain.Block{}
	c := session.DB("blockchain").C("blocks") // use blockchain
	err = c.Find(nil).Limit(15).All(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
