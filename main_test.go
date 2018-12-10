package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/golang-blockchain/blockchain"
	"github.com/stretchr/testify/require"
)

func SendRequest(t *testing.T) (*http.Response, error) {
	reader := strings.NewReader("")
	request, err := http.NewRequest("GET", "http://localhost:8080/", reader)
	if err != nil {
		t.Error(err)
	}
	client := &http.Client{}
	return client.Do(request)
}

func TestMain(t *testing.T) {
	r, err := SendRequest(t)
	if err != nil {
		t.Error("Please run a peer")
		t.FailNow()
	}

	var chain []blockchain.Block
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&chain); err != nil {
		t.Error("Wrong request")
	}

	log.Print(r.Body)
	log.Print(chain)
	require.Equal(t, 1, len(chain))
	require.Equal(t, 0, chain[0].BPM)
}
