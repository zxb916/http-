package server

import (
	"bufio"
	"fmt"
	"glt/com.weiai.www/common/config"
	"glt/com.weiai.www/common/log"
	"glt/com.weiai.www/common/util"
	ants "glt/com.weiai.www/threadpool"
	"io"
	"sync"
	"time"
)

var messageChannel = make(chan string, 10000)

func readData(path string) {
	f := util.ReadFile(path)
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			messageChannel <- "end"
			break
		}
		str := string(line)
		//fmt.Println(str)
		messageChannel <- str
	}

}

func mock() {
	ticker := time.NewTicker(time.Second * 10)
	tk := sync.WaitGroup{}

	tk.Add(config.AppConfig.Server.RequestCount)
	go func() {
		for t := range ticker.C {
			fmt.Printf("send at %s\n", t)
			messageChannel <- config.AppConfig.Files.Url
			tk.Done()
		}
	}()
	tk.Wait()
	ticker.Stop()
	messageChannel <- "end"
}

func sendRequest() {

	// Use the common pool.
	var wg sync.WaitGroup

	p, _ := ants.NewPoolWithFunc(config.AppConfig.Server.ThreadCount, func(line interface{}) {
		send(line)
		wg.Done()
	})
	for {
		url := <-messageChannel
		if url == "end" {
			break
		}
		wg.Add(1)
		_ = p.Invoke(url)
	}
	wg.Wait()
	resultChannel <- "end"
}

func send(i interface{}) {
	url := i.(string)
	time.Sleep(10 * time.Millisecond)
	resp, err := goClient.Get(url)
	if err != nil {
		// handle error
		log.Error(err)
	}
	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	// handle error
	//	log.Error(err)
	//}
	log.Info(resp.StatusCode)
	//fmt.Println(body)

}
