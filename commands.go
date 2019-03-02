package main

import (
	"fmt"
	"time"
	"bytes"
)

//实现具体的命令

func (cli *CLI) AddBlock(txs []*Transaction) {
	cli.bc.AddBlock(txs)
	fmt.Printf("添加区块成功!\n")
}

func (cli *CLI) PrintChain() {

	it := cli.bc.NewIterator()

	for {
		block := it.Next()
		fmt.Printf("++++++++++++++++++++++++++++++++\n")

		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PrevBlockHash)
		fmt.Printf("MerKleRoot : %x\n", block.MerKleRoot)

		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp : %s\n", timeFormat)

		fmt.Printf("Difficulity : %d\n", block.Difficulity)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Transactions[0].TXInputs[0].Address)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid: %v\n", pow.IsValid())

		if bytes.Equal(block.PrevBlockHash, []byte{}) {
			fmt.Printf("区块链遍历结束!\n")
			break
		}
	}
}

func (cli *CLI) Send(from, to string, amount float64, miner string, data string) {

	//1. 创建挖矿交易
	coinbase := NewCoinbaseTx(miner, data)

	txs := []*Transaction{coinbase}

	//2. 创建普通交易
	tx := NewTransaction(from, to, amount, cli.bc)

	if tx != nil {
		txs = append(txs, tx)
	} else {
		fmt.Printf("发现无效交易，过滤!\n")
	}

	//3. 添加到区块
	cli.bc.AddBlock(txs)

	fmt.Printf("挖矿成功!")
}
