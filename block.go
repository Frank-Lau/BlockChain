package main

import "crypto/sha256"

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