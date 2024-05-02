package main

import (
	"Icetea/config"
	"fmt"
	"math/big"
	"os"

	_ "github.com/ethereum/go-ethereum"
	"github.com/joho/godotenv"
)

const (
	BSC_RPC          = "https://bsc-mainnet.nodereal.io/v1/5efb15bec1064001aa27a06c8cedd725"
	CONTRACT_ADDRESS = "CONTRACT_ADDRESS"
	StartBlockNumber = int64(20981361)
)

func init() {
	// Load .env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	contractAddress := os.Getenv(CONTRACT_ADDRESS)
	// Connect Db
	// DSN := os.Getenv("DSN")
	// config.ConnectDb(DSN)

	// err := config.CreateTable()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// err = config.CreateIndexes()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	err := config.IndexingEvent(BSC_RPC, contractAddress, big.NewInt(StartBlockNumber))
	if err != nil {
		fmt.Println(err)
	}

}
