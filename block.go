package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

/*
区块高度
时间戳
数据
哈希
上一个区块的哈希
挖矿的难度值，
计数器，密码学术语
*/

// Block represents each 'item' in the blockchain
// Block就是区块链中的“块”，其中的每一个术语在我的博客中都由详细的解释
type Block struct {
	Index      int
	Timestamp  string
	BPM        int
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
}

/*
	To achieve the
	为了实现区块链的持久化，我们使用Redis数据库来保存
	为了可以把区块链放进Redis，我们需要对区块进行序列化

*/

// 将 Block 序列化为一个字节数组
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 将字节数组反序列化为一个 Block
func DeserializeBlock(d []byte) Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return block
}
