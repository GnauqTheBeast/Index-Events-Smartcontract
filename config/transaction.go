package config

import (
	"context"
	"fmt"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	BlockRange = 50000
)

var (
	LogRequestCreated  = crypto.Keccak256Hash([]byte("RequestCreated(address,uint256,uint256)")).Hex()
	LogResponseCreated = crypto.Keccak256Hash([]byte("ResponseCreated(address,uint256,uint256[])")).Hex()
)

func IndexingEvent(BSC_RPC, ContractAddress string, StartBlockNumber *big.Int) error {
	client, err := ethclient.Dial(BSC_RPC)
	if err != nil {
		return err
	}

	currentBlock := StartBlockNumber

	contractAddress := common.HexToAddress(ContractAddress)

	for {
		// If block not exist, sleep
		err = Delay(client, currentBlock)
		if err != nil {
			return err
		}

		query := ethereum.FilterQuery{
			FromBlock: currentBlock,
			ToBlock:   currentBlock.Add(currentBlock, big.NewInt(1000)),
			Addresses: []common.Address{
				contractAddress,
			},
		}
		fmt.Println(currentBlock)

		err = IterateLogs(client, query)
		if err != nil {
			return err
		}
	}
}

func Delay(client *ethclient.Client, blockNumber *big.Int) error {
	latestBlock, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	if timeDelay := blockNumber.Int64() - latestBlock.Number.Int64(); timeDelay > 0 {
		time.Sleep(time.Duration(int(timeDelay)*3+1) * time.Second)
	}

	return nil
}

func IterateLogs(client *ethclient.Client, query ethereum.FilterQuery) error {
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}

	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case LogRequestCreated:
			err = SaveRequestToDb(vLog)
			if err != nil {
				return err
			}

		case LogResponseCreated:
			SaveResponseToDb(vLog)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SaveRequestToDb(vLog types.Log) error {
	user := common.HexToAddress(vLog.Topics[1].Hex()).String()
	requestId := common.HexToAddress(vLog.Topics[2].Hex()).String()

	request := &Request{
		BlockNumber: int64(vLog.BlockNumber),
		TxHash:      vLog.TxHash.String(),
		User:        user,
		RequestId:   requestId,
		Amount:      int(new(big.Int).SetBytes(vLog.Data).Int64()),
		TxIndex:     int(vLog.Index),
	}

	// Insert to Db
	fmt.Println("Insert to Db Request....")
	_, err := db.NewInsert().Model(request).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Insert Successful")

	return nil
}

func SaveResponseToDb(vLog types.Log) error {
	user := common.HexToAddress(vLog.Topics[1].Hex()).String()
	requestId := common.HexToAddress(vLog.Topics[2].Hex()).String()
	prizeIds := make([]int, 0)
	for i := 0; i < len(vLog.Data); i += 32 {
		prizeIds = append(prizeIds, int(new(big.Int).SetBytes(vLog.Data[i:i+32]).Int64()))
	}

	response := &Response{
		BlockNumber: int64(vLog.BlockNumber),
		TxHash:      vLog.TxHash.String(),
		User:        user,
		RequestId:   requestId,
		PrizeIds:    prizeIds,
		TxIndex:     int(vLog.Index),
	}

	// Insert to Db
	fmt.Println("Insert to Db Respone....")
	_, err := db.NewInsert().Model(response).Exec(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Insert Successful")

	return nil
}
