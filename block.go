package main

import (
	"time"
	"encoding/gob"
	"bytes"
	"fmt"
	"log"
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

//序列化, 将区块转换成字节流
func (block *Block) Serialize() []byte {

	var buffer bytes.Buffer

	//定义编码器
	encoder := gob.NewEncoder(&buffer)

	//编码器对结构进行编码，一定要进行校验
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

func Deserialize(data []byte) *Block {

	fmt.Printf("解码传入的数据: %x\n", data)

	var block Block

	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}
