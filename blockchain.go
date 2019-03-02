package main

import (
	"./bolt"
	"log"
	"fmt"
	"os"
)

//使用bolt进行改写，需要两个字段：
//1. bolt数据库的句柄
//2. 最后一个区块的哈希值
type BlockChain struct {
	db   *bolt.DB //句柄
	tail []byte   //最后一个区块的哈希值
}

const blockChainName = "blockChain.db"
const blockBucketName = "blockBucket"
const lastHashKey = "lastHashKey"

//实现创建区块链的方法
func NewBlockChain(miner string) *BlockChain {

	//功能分析：
	//1. 获得数据库的句柄，打开数据库，读写数据

	db, err := bolt.Open(blockChainName, 0600, nil)
	//向数据库中写入数据
	//从数据库中读取数据

	if err != nil {
		log.Panic(err)
	}

	//defer db.Close()

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
			//创世块中只有一个挖矿交易，只有Coinbase
			coinbase := NewCoinbaseTx(miner, genesisInfo)
			genesisBlock := NewBlock([]*Transaction{coinbase}, []byte{})

			b.Put(genesisBlock.Hash, genesisBlock.Serialize() /*将区块序列化，转成字节流*/)
			b.Put([]byte(lastHashKey), genesisBlock.Hash)

			//为了测试，我们把写入的数据读取出来，如果没问题，注释掉这段代码
			//blockInfo := b.Get(genesisBlock.Hash)
			//block := Deserialize(blockInfo)
			//fmt.Printf("解码后的block数据:%s\n", block)

			tail = genesisBlock.Hash
		} else {
			tail = b.Get([]byte(lastHashKey))
		}

		return nil
	})

	return &BlockChain{db, tail}
}

//添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//1. 创建一个区块
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))

		if b == nil {
			fmt.Printf("bucket不存在，请检查!\n")
			os.Exit(1)
		}

		block := NewBlock(txs, bc.tail)
		b.Put(block.Hash, block.Serialize() /*将区块序列化，转成字节流*/)
		b.Put([]byte(lastHashKey), block.Hash)

		bc.tail = block.Hash

		return nil
	})
}

//定义一个区块链的迭代器，包含db，current
type BlockChainIterator struct {
	db      *bolt.DB
	current []byte //当前所指向区块的哈希值
}

//创建迭代器，使用bc进行初始化

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{bc.db, bc.tail}
}

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

//实现思路：
//

func (bc *BlockChain) FindMyUtoxs(address string) []TXOutput {
	fmt.Printf("FindMyUtoxs\n")
	var UTXOs []TXOutput //返回的结构

	it := bc.NewIterator()

	//这是标识已经消耗过的utxo的结构，key是交易id，value是这个id里面的output索引的数组
	spentUTXOs := make(map[string][]int64)

	//1. 遍历账本
	for {

		block := it.Next()

		//2. 遍历交易
		for _, tx := range block.Transactions {
			//遍历交易输入:inputs

			for _, input := range tx.TXInputs {
				if input.Address == address {
					fmt.Printf("找到了消耗过的output! index : %d\n", input.Index)
					key := string(input.TXID)
					spentUTXOs[key] = append(spentUTXOs[key], input.Index)
					//spentUTXOs[0x222] = []int64{0}
					//spentUTXOs[0x333] = []int64{0}  //中间状态
					//spentUTXOs[0x333] = []int64{0, 1}
				}
			}

			key := string(tx.TXid)
			indexes /*[]int64{0,1}*/ := spentUTXOs[key]

		OUTPUT:
		//3. 遍历output
			for i, output := range tx.TXOutputs {

				if len(indexes) != 0 {
					fmt.Printf("当前这笔交易中有被消耗过的output!\n")
					for _, j /*0, 1*/ := range indexes {
						if int64(i) == j {
							fmt.Printf("i == j, 当前的output已经被消耗过了，跳过不统计!\n")
							continue OUTPUT
						}
					}
				}

				//4. 找到属于我的所有output
				if address == output.Address {
					fmt.Printf("找到了属于 %s 的output, i : %d\n", address, i)
					UTXOs = append(UTXOs, output)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("遍历区块链结束!\n")
			break
		}
	}

	return UTXOs
}

func (bc *BlockChain) GetBalance(address string) {
	utxos := bc.FindMyUtoxs(address)

	var total = 0.0

	for _, utxo := range utxos {
		total += utxo.Value //10, 3, 1
	}

	fmt.Printf("%s 的余额为: %f\n", address, total)
}

//1. 遍历账本，找到属于付款人的合适的金额，把这个outputs找到
//utxos, resValue = bc.FindNeedUtxos(from, amount)
func (bc *BlockChain) FindNeedUtxos(from string, amount float64) (map[string][]int64, float64) {

	needUtxos := make(map[string][]int64) //标识能用的utxo, //返回的结构
	var resValue float64                  //统计的金额

	it := bc.NewIterator()

	//这是标识已经消耗过的utxo的结构，key是交易id，value是这个id里面的output索引的数组
	spentUTXOs := make(map[string][]int64)

	//1. 遍历账本
	for {

		block := it.Next()

		//2. 遍历交易
		for _, tx := range block.Transactions {
			//遍历交易输入:inputs

			for _, input := range tx.TXInputs {
				if input.Address == from {
					fmt.Printf("找到了消耗过的output! index : %d\n", input.Index)
					key := string(input.TXID)
					spentUTXOs[key] = append(spentUTXOs[key], input.Index)
					//spentUTXOs[0x222] = []int64{0}
					//spentUTXOs[0x333] = []int64{0}  //中间状态
					//spentUTXOs[0x333] = []int64{0, 1}
				}
			}

			key := string(tx.TXid)
			indexes /*[]int64{0,1}*/ := spentUTXOs[key]

		OUTPUT:
		//3. 遍历output
			for i, output := range tx.TXOutputs {

				if len(indexes) != 0 {
					fmt.Printf("当前这笔交易中有被消耗过的output!\n")
					for _, j /*0, 1*/ := range indexes {
						if int64(i) == j {
							fmt.Printf("i == j, 当前的output已经被消耗过了，跳过不统计!\n")
							continue OUTPUT
						}
					}
				}

				//4. 找到属于我的所有output
				if from == output.Address {
					fmt.Printf("找到了属于 %s 的output, i : %d\n", from, i)
					//UTXOs = append(UTXOs, output)
					//在这里实现控制逻辑
					//找到符合条件的output
					//1. 添加到返回结构中needUtxos
					needUtxos[key] = append(needUtxos[key], int64(i))
					resValue += output.Value

					//2. 判断一下金额是否足够
					if resValue >= amount {
						//a. 足够， 直接返回
						return needUtxos, resValue
					}
					//b. 不足， 继续遍历
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("遍历区块链结束!\n")
			break
		}
	}

	return needUtxos, resValue
}
