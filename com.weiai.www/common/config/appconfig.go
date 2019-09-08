package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"glt/com.weiai.www/common/log"
	"os"
	"regexp"
	"time"
)

const configFileSizeLimit = 10 << 20

var AppConfig *config

//var once sync.Once

//有了`json:network`这种注释，后面json解析就可以把相应的数据塞到对应的结构里面来
type config struct {
	Server serverConfig `json:server`
	Files  fileConfig   `json:files`
}

type serverConfig struct {
	ThreadCount  int           `json:ThreadCount`
	MaxIdleConn  int           `json:MaxIdleConn"`
	RequestCount int           `json:"RequestCount"`
	Debug        bool          `json:"Debug"`
	Timeout      time.Duration `json:timeout`
}

type fileConfig struct {
	Paths    string            `json:paths`
	Fields   map[string]string `json:fields`
	DeadTime time.Duration
	Url      string `json:"url"`
}

//func init() {
//	loadConfig("../../conf.json")
//}
func LoadConfig(path string) {
	config_file, err := os.Open(path)
	if err != nil {
		log.Error("Failed to open config file '%s': %s\n", path, err)
		return
	}

	fi, _ := config_file.Stat()
	if size := fi.Size(); size > (configFileSizeLimit) {
		log.Error("config file (%q) size exceeds reasonable limit (%d) - aborting", path, size)
		return // REVU: shouldn't this return an error, then?
	}

	if fi.Size() == 0 {
		log.Error("config file (%q) is empty, skipping", path)
		return
	}

	buffer := make([]byte, fi.Size())
	_, err = config_file.Read(buffer)
	if err != nil {
		log.Error("\n %s\n", buffer)
	}

	buffer, err = stripComments(buffer) //去掉注释
	if err != nil {
		log.Error("Failed to strip comments from json: %s\n", err)
		return
	}

	buffer = []byte(os.ExpandEnv(string(buffer))) //特殊

	err = json.Unmarshal(buffer, &AppConfig) //解析json格式数据
	if err != nil {
		log.Error("Failed unmarshalling json: %s\n", err)
		return
	}
	fmt.Printf("111111111 %s \n", AppConfig.Server.ThreadCount)

	if AppConfig.Files.DeadTime == 0 {
		AppConfig.Files.DeadTime = 60
	}

	return
}

//去掉注释
func stripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0) // Windows
	lines := bytes.Split(data, []byte("\n"))                //split to muli lines
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}

	return bytes.Join(filtered, []byte("\n")), nil
}
