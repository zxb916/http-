package server

import (
	"fmt"
	"glt/com.weiai.www/common/config"
	"glt/com.weiai.www/common/log"
	"time"
)

var resultChannel = make(chan string, 10000)

func Start() {
	//初始化httpClient
	initClient()
	if config.AppConfig.Server.Debug {
		go mock()
	} else {
		go readData(config.AppConfig.Files.Paths)
	}
	run()
}

func run() {
	start := time.Now()
	go sendRequest()
	result := <-resultChannel
	if result == "end" {
		log.Info("complete task ")
	}
	elapsed := time.Since(start)
	fmt.Println("总共耗时:", elapsed)
}
