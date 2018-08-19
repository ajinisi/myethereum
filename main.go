package main

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
)

var mutex = &sync.Mutex{}
var bc *Blockchain

func main() {
	// 完成环境变量（:8080端口）的加载，以便通过浏览器进行访问
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// 防止数据库没有区块链，预先生成传世块
	bc = NewBlockchain()

	log.Fatal(run())

}

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}
