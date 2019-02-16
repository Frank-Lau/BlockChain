package main

import (
	"fmt"
	"crypto/sha256"
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

	bc.AddBlcok("老王来了")

	for i,block := range bc.Blocks{
		fmt.Printf("+++++++++++++++++++++++%d+++++++++++++++++++++++++++\n",i)
		fmt.Printf("PrevBlockHash : %x\n", block.PrevBlockHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Data)
	}


}

