package main

/*
 * ToTo - Go language proxy server
 * Song Hyeon Sik <blogcin@naver.com> 2015
 */

import (
	"fmt"
	"net"
	"strconv"
)

const (
	bufferLength = 8192
)

type ProxyServer struct {
	port string
}

func (ps *ProxyServer) askPort() int{
	port := 0

	fmt.Print("Port : ")
	fmt.Scanf("%d", &port)

	return port
}

func (ps *ProxyServer) init(port int) net.Listener {
	ps.port = ":"
	ps.port += strconv.Itoa(port)

	server, err := net.Listen("tcp", ps.port)

	if server == nil {
		panic("init: port listening error : " + err.Error())
	}

	return server
}

func (ps *ProxyServer) acceptClient(server net.Listener) chan net.Conn{
	channel := make(chan net.Conn)

	go func() {
		for {
			client, err := server.Accept()
			if client == nil {
				fmt.Printf("ps: acceptClient: Couldn't accept : " + err.Error())
				continue
			}
			channel <- client
		}

	}()
	return channel
}

func (ps *ProxyServer) connectHost(client net.Conn) {
	ps.getHeader(client)

	// Get Header
	//fmt.Println(string(buffer[:]))
	//conn, _ := net.Dial("tcp", "127.0.0.1" + ps.port)

}

func (ps *ProxyServer) getHeader(client net.Conn) {
	buffer := make([]byte, bufferLength)
	client.Read(buffer)

	fmt.Println(buffer)
}

/*
func byteArrtoStr(byteArray []byte) string{
	return string(byteArray[:])
}
*/
func main() {
	proxyServer := &ProxyServer{}

	port := proxyServer.askPort()

	server := proxyServer.init(port)
	defer server.Close()

	connections := proxyServer.acceptClient(server)

	for {
		go proxyServer.connectHost(<-connections)
	}
}
