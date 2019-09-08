package server

import (
	"crypto/tls"
	"glt/com.weiai.www/common/config"
	"net/http"
)

var goClient *http.Client

func initClient() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:    config.AppConfig.Server.MaxIdleConn,
	}
	goClient = &http.Client{Transport: transport}
}
