package utils

import (
	"log"

	"github.com/btcsuite/btcd/rpcclient"
)

func Bitcoind() *rpcclient.Client {

	// Connect to a running Bitcoin Core node via RPC
	connCfg := &rpcclient.ConnConfig{
		Host:         GetEnv("BTC_HOST"),
		User:         GetEnv("BTC_USER"),
		Pass:         GetEnv("BTC_PASS"),
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal("Error connecting to bitcoind:", err)
	}
	// Get the current block count
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Println("Error getting block count:", err)

	}

	log.Println("Current block count:", blockCount)
	return client
}
