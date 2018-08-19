package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

// Message takes incoming JSON payload for writing heart rate
type Message struct {
	BPM int
}

// web server
func run() error {
	mux := makeMuxRouter()
	// 从我们前面创建的.env文件中提取:8081
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
	//muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", bc.handleWriteBlock).Methods("POST")
	return muxRouter
}

// write blockchain when we receive an http request
// 当我们收到一个http请求时，返回整个区块链
// func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
// 	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	io.WriteString(w, string(bytes))
// }

// takes JSON payload as an input for heart rate (BPM)
// 接收JSON数据作为BPM的输入
func (bc *Blockchain) handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message

	// decoder := json.NewDecoder(r.Body)
	// if err := decoder.Decode(&m); err != nil {
	// 	respondWithJSON(w, r, http.StatusBadRequest, r.Body)
	// 	return
	// }
	// defer r.Body.Close()

	r.ParseForm() // 解析参数，默认是不会解析的
	// 在控制台上输出信息
	fmt.Println("Form: ", r.Form)
	fmt.Println("Path: ", r.URL.Path)
	// 接收到的数据处理
	m.BPM, _ = strconv.Atoi(r.Form["BPM"][0])

	// 首先获取最后一个块用于生成新块
	c := bc.c
	tip := bc.tip
	fmt.Println("Get lastBlock Hash: ", tip)
	lastBlock, err := redis.Bytes(c.Do("GET", tip))
	if err != nil {
		fmt.Println(err)
	}
	// ensure atomicity when creating new block
	// 写入新的区块之前，我们需要锁定互斥量
	mutex.Lock()
	newBlock := generateBlock(DeserializeBlock(lastBlock), m.BPM)
	mutex.Unlock()

	// 当确认新区块合法后，如果合法则把它加入数据库
	if isBlockValid(newBlock, DeserializeBlock(lastBlock)) {
		_, err := c.Do("MSET", newBlock.Hash, newBlock.Serialize(), "l", newBlock.Hash)
		if err != nil {
			fmt.Println("redis set failed:", err)
		}
		// 更新最后一块的哈希
		bc.tip = newBlock.Hash
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

// 一旦API调用过程中出现错误就能以JSON格式返回错误信息
func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	// w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}
