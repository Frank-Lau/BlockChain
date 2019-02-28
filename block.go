package main

import (
	"time"
)

const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Block struct {
	Version       uint64 //版本号
	PrevBlockHash []byte //前区块哈希
	MerKleRoot    []byte //先填写为空,后续v4版本补充
	TimeStamp     uint64 //从1970.1.1自己秒数
	Difficulity   uint64 //挖矿难度值,v2时使用
	Nonce         uint64 //随机数,挖矿找的就是他
	Data          []byte //数据，目前使用字节流，v4开始使用交易代替
	Hash          []byte //当前区块哈希,区块中本不存在的字段,为了方便我们添加进来
}

//创建区块，对Block的每一个字段填充数据即可
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:       00,
		PrevBlockHash: prevBlockHash,
		MerKleRoot:    []byte{},
		TimeStamp:     uint64(time.Now().Unix()),
		Difficulity:   Bits, //随便写的,v2在调整
		//Nonce:         10, //同difficulity
		Data:          []byte(data),
		Hash:          []byte{}, //先填充为空，后续会填充数据
	}
	//调用setHash生成哈希值
	//block.setHash()
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

//v1版本中的此函数已在pow函数中实现,所以注释掉
/*
//实现setHash函数,我们实现一个简单的函数,来计算哈希子,没有随机数,没有难度值
func (block *Block) setHash() {
	*/
/*var data []byte

	//uintToByte将数字转化为[]byte{},在utils中实现
	data = append(data, uintToByte(block.Version)...)
	data = append(data, block.PrevBlockHash...)
	data = append(data, block.MerKleRoot...)
	data = append(data, uintToByte(block.TimeStamp)...)
	data = append(data, uintToByte(block.Difficulity)...)
	data = append(data, uintToByte(block.Nonce)...)
	data = append(data, block.Data...)*//*



	//使用byte.join改写setHash
	tmp := [][]byte{
		uintToByte(block.Version),
		block.PrevBlockHash,
		block.MerKleRoot,
		uintToByte(block.TimeStamp),
		uintToByte(block.Difficulity),
		uintToByte(block.Nonce),
		block.Data,
	}
	//传入一个二位切片,以一个以为切片进行分割,并返回一个一维切片
	data  :=bytes.Join(tmp,[]byte{})

	hash */
/*[32]byte*//*
 := sha256.Sum256(data)
	block.Hash = hash[:]
}
*/
