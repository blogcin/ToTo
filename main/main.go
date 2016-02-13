package main

import (
	"github.com/blogcin/ToTo/ToTo"
)
/*
 * ToTo - Go language proxy server
 * Song Hyeon Sik <blogcin@naver.com> 2016
 */


func main() {
	ToTo := &ToTo{}

	port := ToTo.askPort()

	server := ToTo.init(port)
	defer server.Close()

	connections := ToTo.acceptClient(server)

	for {
		go ToTo.connectHost(<-connections)
	}
}
