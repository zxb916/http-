// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"glt/com.weiai.www/common/config"
	"glt/com.weiai.www/server"
)

func main() {

	//加载配置文件
	config.LoadConfig("D:/DevelopSoft/goidea/workspace/goclient/com.weiai.www/conf.json")
	server.Start()

}
