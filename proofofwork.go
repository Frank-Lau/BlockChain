package main

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

//1.定义一个工作量证明的结构ProofOfWork
//	a.block
//	b.目标值
//2.提供创建POW函数
//	NewProofOfWork(参数)
//3.提供不断计算hash的	函数
//	run()
//4.提供一个校验函数

//POW定义
type ProofOfWork struct {
	block *Block

	//采用big.int来存储哈希值,因为哈希值太大,内置的方法Cmp:比较方法
	//SetBytes:把Bytes转化成big.int
	//[]byte("0x00000919011eeb8fbdf0c476d8510b8e1e632eba7b584ac04c11ad20cbbdd394")

	//SetString:baString转化成big.int
	//"0x00000919011eeb8fbdf0c476d8510b8e1e632eba7b584ac04c11ad20cbbdd394"
	target *big.Int //系统提供的,固定的,大约两周调整一次
}

//2.提供创建POW的函数
const Bits  =16
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}


	/*//写难度值,难度值应该是推导出来的,但是我们为了简化,把难度值先写成固定的,一切完成之后,再去推到
	// 0000100000000000000000000000000000000000000000000000000000000000
	//16制格式的字符串转化成big.int类型
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	var bigIntTmp big.Int
	bigIntTmp.SetString(targetStr, 16)*/


	//使用程序来推导难度值,推导前导另为三个的难度值
	//
	//0001000000000000000000000000000000000000000000000000000000000000
	//初始化
	//0001000000000000000000000000000000000000000000000000000000000001
	//将1向左移动256位
	//1 0000000000000000000000000000000000000000000000000000000000000000
	//向右移动四次,一个16禁止位代表4个2进制为(发:1111)
	//0 0001000000000000000000000000000000000000000000000000000000000000

	bigIntTmp := big.NewInt(1)
	//bigIntTmp.Lsh(bigIntTmp,256)
	//bigIntTmp.Rsh(bigIntTmp,16)
	bigIntTmp.Lsh(bigIntTmp,256-Bits)


	pow.target = bigIntTmp


	return &pow
}

//3.创建不断计算哈希的run函数.为了获取挖矿的随机数,并且返回区块的哈希值
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	//1.获取block数据
	//2.拼接nonce
	//3.sha256
	//4.与难度值比较
	//	a.哈希值大于难度值,nonce++
	//	b.哈希值小于难度值,挖矿成功退出
	var nonce uint64
	var hash [32]byte
	for ; ; {
		fmt.Printf("%x\r",hash)
		//data = block +nonce
		hash = sha256.Sum256(pow.prepareData(nonce))
		//将hash类型转化成big.int,然后与pow.target进行比较,需要引入bigIntTmp局部变量
		var bigIntTmp big.Int
		bigIntTmp.SetBytes(hash[:])
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//
		//func (x *Int) Cmp(y *Int) (r int) {
		//   x              y
		if bigIntTmp.Cmp(pow.target) == -1 {
			//此时x<y,挖矿成功
			fmt.Printf("挖矿成功!nonce:%d,哈希值为:%x\n", nonce, hash)
			break
		} else {
			nonce++
		}

	}
	return hash[:],nonce

}

//定义准备数据函数prepare
func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	block := pow.block
	tmp := [][]byte{
		uintToByte(block.Version),
		block.PrevBlockHash,
		block.MerKleRoot,
		uintToByte(block.TimeStamp),
		uintToByte(block.Difficulity),
		block.Data,
		uintToByte(nonce),
	}
	data := bytes.Join(tmp, []byte{})
	return data
}

//检验挖矿是否有效函数IsValid
func (pow *ProofOfWork)IsValid()bool{
	//在校验的时候,block的数据是完整的,我们要做的就是检查一下,hash,block数据和Nonce是否满足难度值要求

	//获取区块数据
	//拼接nonce
	//做sha256

	data :=pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)

	var tmp big.Int
	tmp.SetBytes(hash[:])

	//if tmp.Cmp(pow.target)==-1{
	//	return true
	//}else {
	//	return false
	//}

	return tmp.Cmp(pow.target)==-1

}


