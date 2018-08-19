package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

/*
	为了持久化，我们区块链的结构体定义需要修改如下
	tip 这个词本身有事物尖端或尾部的意思，这里指的是存储最后一个块的哈希
	在链的末端可能出现短暂分叉的情况，所以选择 tip 其实也就是选择了哪条链
	db 存储数据库连接
*/

// Blockchain is a series of validated Blocks
// 区块链就是被验证的区块的一串序列

/*
	我们不再将区块链储存在一个数组里
	我们现在所能操作的仅仅是区块链的最后一个区块的哈希
	同时为了方便，也存储一个数据库的连接
*/
// var Blockchain []Block

type Blockchain struct {
	tip string
	c   redis.Conn
}

// 用来遍历区块链
type BlockchainIterator struct {
	CurrentHash string
	c           redis.Conn
}

func NewBlockchain() *Blockchain {
	var tip string

	// 打开redis
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// handle error
	}
	// 不再关闭连接，因为后面还会用到
	// defer c.Close()

	if _, err := c.Do("AUTH", "123456"); err != nil {
		// handle error
	}

	is_key_exit, err := redis.Bool(c.Do("EXISTS", "l"))
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("exists or not: %v \n", is_key_exit)
	}
	// 如果数据库中不存在区块链就创建一个，否则直接读取最后一个块的哈希
	if is_key_exit == false {
		fmt.Println("No existing blockchain found. Creating a new one...")
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), "", difficulty, ""}

		_, err := c.Do("MSET", genesisBlock.Hash, genesisBlock.Serialize(), "l", genesisBlock.Hash)

		if err != nil {
			fmt.Println("redis set failed:", err)
		}
		tip = genesisBlock.Hash

	} else {
		var err error
		tip, err = redis.String(c.Do("GET", "l"))
		if err != nil {
			fmt.Println("redis get failed:", err)
		} else {
			fmt.Printf("Get mykey: %v \n", tip)
		}
	}

	bc := &Blockchain{
		tip,
		c,
	}
	return bc
}

/*

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// 返回链中的下一个块
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

*/
