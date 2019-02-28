package main

import (
	"fmt"



)

//1. 定义结构（区块头的字段比正常的少）
//>1. 前区块哈希
//>2. 当前区块哈希
//>3. 数据

//2. 创建区块
//3. 生成哈希
//4. 引入区块链
//5. 添加区块
//6. 重构代码

func main() {
	fmt.Printf("helloworld\n")

	//block := NewBlock(genesisInfo, []byte{0x0000000000000000})
	bc := NewBlockChain()
	defer bc.db.Close()
	cli := CLI{bc}
	cli.Run()
/*	bc.AddBlock("hello itcast!!!")

	it := bc.NewIterator()

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
		fmt.Printf("Data : %s\n", block.Data)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid: %v\n", pow.IsValid())

		if bytes.Equal(block.PrevBlockHash, []byte{}) {
			fmt.Printf("区块链遍历结束!\n")
			break
		}
	}*/
}