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

const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Block struct {
	PrevBlockHash []byte //前区块哈希
	Hash          []byte //当前区块哈希
	Data          []byte //数据，目前使用字节流，v4开始使用交易代替
}

//创建区块，对Block的每一个字段填充数据即可
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{}, //先填充为空，后续会填充数据
		Data:          []byte(data),
	}
	//调用setHash生成哈希值
	block.setHash()

	return &block
}

//实现setHash函数,我们实现一个简单的函数,来计算哈希子,没有随机数,没有难度值
func (block *Block)setHash(){
	var data []byte
	data = append(data,block.Data...)
	data = append(data,block.PrevBlockHash...)

	hash /*[32]byte*/ := sha256.Sum256(data)
	block.Hash  = hash[:]
}

//区块链的定义及遍历打印

//创建区块链,使用Block切片模拟
type BlockChain struct {
	Blocks []*Block
}

//实现创建区块链的方法
func NewBlockChain()*BlockChain{
	//在创建的时候添加一个区块,创世块
	genesisBlock := NewBlock(genesisInfo,[]byte{0x0000000000000000})

	bc :=BlockChain{Blocks:[]*Block{genesisBlock}}
	return &bc
}


//添加区块
func (bc *BlockChain)AddBlcok(data string){
	//1.创建一个区块

	//bc.Blocks的最后一个区块的哈希就是当前新区块的PrevBlockHash
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	prevHash :=lastBlock.Hash

	block := NewBlock(data,prevHash)


	//2.添加到bc.Blocks数组中
	bc.Blocks =append(bc.Blocks,block)
}

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

