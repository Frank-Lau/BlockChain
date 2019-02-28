package main

import (
	"github.com/bolt"
	"log"
	"fmt"
	"os"
)

//区块链的定义及遍历打印

////创建区块链,使用Block切片模拟
//type BlockChain struct {
//	Blocks []*Block
//}

//使用bolt数据库进行持久化存储
type BlockChain struct {
	db   *bolt.DB //句柄
	tail []byte   //最后一个区块的哈希值
}

const blockChainName = "blockChain.db"
const blockBucketName = "blockBucket"
const lastHashKey = "lastHashKey"

//实现创建区块链的方法
func NewBlockChain() *BlockChain {
	//功能分析：
	//1. 获得数据库的句柄，打开数据库，读写数据

	db, err := bolt.Open(blockChainName, 0600, nil)
	//向数据库中写入数据
	//从数据库中读取数据

	if err != nil {
		log.Panic(err)
	}

	//defer db.Close()不能关闭,因为后面要使用

	var tail []byte

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))

		if b == nil {
			//如果b1为空，说明名字为"buckeName1"这个桶不存在，我们需要创建之
			fmt.Printf("bucket不存在，准备创建!\n")
			b, err = tx.CreateBucket([]byte(blockBucketName))

			if err != nil {
				log.Panic(err)
			}

			//抽屉准备完毕，开始添加创世块
			genesisBlock := NewBlock(genesisInfo, []byte{})
			b.Put(genesisBlock.Hash, genesisBlock.Serialize() /*将区块序列化，转成字节流*/)
			b.Put([]byte(lastHashKey), genesisBlock.Hash)

			tail = genesisBlock.Hash
		} else {
			tail = b.Get([]byte(lastHashKey))
		}

		return nil
	})

	return &BlockChain{db, tail}
}

//添加区块
func (bc *BlockChain) AddBlock(data string) {
	//1. 创建一个区块
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))

		if b == nil {
			fmt.Printf("bucket不存在，请检查!\n")
			log.Panic()
		}

		block := NewBlock(data, bc.tail)
		b.Put(block.Hash, block.Serialize() /*将区块序列化，转成字节流*/)
		b.Put([]byte(lastHashKey), block.Hash)

		bc.tail = block.Hash

		return nil
	})
}
//定义一个区块链的迭代器，包含db，current,跟数组和切片相反,是从后往前遍历
type BlockChainIterator struct {
	db      *bolt.DB
	current []byte //当前所指向区块的哈希值
}

//创建迭代器，使用bc进行初始化

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{bc.db, bc.tail}
}
//实现next函数,功能一,返回当前区块数据,current前移
func (it *BlockChainIterator) Next() *Block {

	var block Block

	it.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockBucketName))
		if b == nil {
			fmt.Printf("bucket不存在，请检查!\n")
			os.Exit(1)
		}

		//真正的读取数据
		blockInfo /*block的字节流*/ := b.Get(it.current)
		block = *Deserialize(blockInfo)

		it.current = block.PrevBlockHash

		return nil
	})

	return &block
}
