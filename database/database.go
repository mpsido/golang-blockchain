package database

import (
	// "fmt"

	"log"

	mgo "gopkg.in/mgo.v2"

	// "gopkg.in/mgo.v2/bson"
	"github.com/golang-blockchain/blockchain"
)

func WriteBlockchain(inputBlockchain []blockchain.Block) {
	session, err := mgo.Dial("172.17.0.3")
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

	// result := blockchain.Block{}
	// err = c.Find(bson.M{"BPM": 55}).One(&result)
	// if err != nil {
	//         log.Println(err)
	// }
	//        log.Println("Block", result.BPM)

	// fmt.Println("Phone:", result.Phone)
}
