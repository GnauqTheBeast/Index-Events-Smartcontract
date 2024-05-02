package config

import (
	"context"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// type RequestCreated struct {
// 	User      common.Address
// 	RequestId common.Address
// 	Amount    int64
// }

// type ResponseCreated struct {
// 	User      common.Address
// 	RequestId common.Address
// 	PrizeIds  []uint64
// }

var (
	LogRequestCreatedSig  = []byte("RequestCreated(address,uint256,uint256)")
	LogResponseCreatedSig = []byte("ResponseCreated(address,uint256,uint256[])")
)

func IndexingEvent(BSC_RPC, ContractAddress string, StartBlockNumber *big.Int) error {
	client, err := ethclient.Dial(BSC_RPC)
	if err != nil {
		return err
	}

	endBlockNumber := big.NewInt(int64(20981700))

	contractAddress := common.HexToAddress(ContractAddress)
	query := ethereum.FilterQuery{
		FromBlock: StartBlockNumber,
		ToBlock:   endBlockNumber,
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}

	LogRequestCreatedSig := crypto.Keccak256Hash(LogRequestCreatedSig)
	LogResponseCreatedSig := crypto.Keccak256Hash(LogResponseCreatedSig)

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)
		fmt.Printf("TxIndex: %d\n", vLog.TxIndex)
		fmt.Printf("TxHash: %s\n", vLog.TxHash)
		fmt.Printf("TxAddress: %s\n", vLog.Address)

		switch vLog.Topics[0].Hex() {
		case LogRequestCreatedSig.Hex():
			fmt.Printf("Log Name: LogRequestCreated\n")

			User := common.HexToAddress(vLog.Topics[1].Hex())
			RequestId := common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Println(new(big.Int).SetBytes(vLog.Data).Int64())
			fmt.Printf("User: %s\n", User.Hex())
			fmt.Printf("RequestId: %s\n", RequestId.Hex())

		case LogResponseCreatedSig.Hex():
			fmt.Printf("Log Name: LogResponeCreated\n")

			User := common.HexToAddress(vLog.Topics[1].Hex())
			RequestId := common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Println(new(big.Int).SetBytes(vLog.Data).Int64())
			fmt.Printf("User: %s\n", User.Hex())
			fmt.Printf("RequestId: %s\n", RequestId.Hex())
		}
		fmt.Print("\n\n")
	}
	return nil
}
