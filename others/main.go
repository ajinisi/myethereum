package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	// spew 可以帮助我们在 console 中直接查看 struct 和 slice 这两种数据结构
	// 使用spew.Dump 这个函数可以以非常美观和方便阅读的方式将 struct、slice
	// 等数据打印在控制台里，方便我们调试
	"github.com/davecgh/go-spew/spew"
	// Gorilla 的 mux 包非常流行， 我们用它来写 web handler
	"github.com/gorilla/mux"

	// Gotdotenv lets us read from a .env file that we keep in the root of our directory
	// so we don’t have to hardcode things like our http ports.
	// godotenv 可以帮助我们读取项目根目录中的 .env 配置文件，
	// 这样我们就不用将 http port 之类的配置硬编码进代码中了。比如像这样ADDR=8080
	"github.com/joho/godotenv"
)

var mutex = &sync.Mutex{}

// 关于区块链的函数，web 服务的函数“组装”起来
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = *NewGenesisBlock()
		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())

}

// 初始化我们的 web 服务 web server
func run() error {
	mux := makeMuxRouter()
	// 端口号是通过前面提到的 .env 来获得
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	// 对“/”的 GET 请求我们可以查看整个链
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	// “/”的 POST 请求可以创建块
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

// write blockchain when we receive an http request
// GET 请求的 handler 实现
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

// takes JSON payload as an input for heart rate (AGE)

// Message takes incoming JSON payload for writing heart rate
// 先来定义一下 POST 请求的 payload
// 我们的 POST 请求体中可以使用这里定义的 payload，比如：{"AGE":75}
type Message struct {
	AGE int
}

// POST 请求的 handler 实现
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], m.AGE)
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

// POST 请求处理完之后，无论创建块成功与否，我们需要返回客户端一个响应
func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
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

// SHA256 hasing
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.AGE) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, AGE int) Block {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.AGE = AGE
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}
