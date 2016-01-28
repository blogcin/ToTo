package main

/*
 * ToTo - Go language proxy server
 * Song Hyeon Sik <blogcin@naver.com> 2015
 */

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	bufferLength = 15000
	headerLine = 30
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
				fmt.Println("ps: acceptClient: Couldn't accept : ", err.Error())
				continue
			}
			channel <- client
		}

	}()
	return channel
}

func (ps *ProxyServer) connectHost(client net.Conn) {
	HeaderInfo, Datas := ps.getData(client)

	if(Datas[0] == 0 ) || (HeaderInfo == "-1"){
		return
	}
	//buffers := make([]byte, bufferLength)

	requestType, host, protocol, port := ps.parseHttpHeaderMethod(HeaderInfo)

	if(port == -1) {
		return
	}
	fmt.Println("Request Type : ", requestType)
	fmt.Println("HOST : ", host)
	fmt.Println("Protocol : ", protocol)
	fmt.Println("Port : ", port)

	connectionHost, _ := net.Dial("tcp", host + ":"+ strconv.Itoa(port))

	if requestType == "CONNECT" {
		connectionHost.Write([]byte("HTTP/1.1 200 Connection established\n"))
	} else {
		connectionHost.Write(Datas)
		fmt.Println("connectionHost.Write(Datas)")
		for {
			connectionHost.Read(Datas)
			fmt.Println("connectionHost.Read(Datas)")
			client.Write(Datas)
			fmt.Println("client.Write(Datas)")
			client.Read(Datas)
			fmt.Println("client.Read(Datas)")
			connectionHost.Write(Datas)
			fmt.Println("connectionHost.Write(Datas)")

		}
	}
	connectionHost.Close()
	return

	//io.Copy(Datas, connectionHost)
	// Get Header
	//fmt.Println(string(buffer[:]))
	//
}
/*
Get first line of Http header
 */
func (ps *ProxyServer) getData(client net.Conn) (string, []byte){

	buffer := make([]byte, bufferLength)
	client.Read(buffer)


	return ps.splitHeader(buffer)[0], buffer
}

func (ps *ProxyServer) parseHttpHeaderMethod(headerMethod string) (string, string, string, int) {
	var (
		requestType string
		host string
		protocol string
		port int
	)
	// ex: GET http://google.com/ HTTP/1.1

	temp := headerMethod[strings.Index(headerMethod, " ")+1:]
	protocol = temp[strings.Index(temp, " ")+1:]

	url := temp[:strings.Index(temp, " ")]

	i := strings.Index(url, "://")

	if i == -1 {
		fmt.Println("Uncorrect URL")
		return "", "", "", -1
	} else {
		host = url[i+3:len(url)-1]
	}

	i = strings.Index(host, ":")
	if i == -1 {
		port = 80
	} else {
		port, _ = strconv.Atoi(host[i+1:])
		host = host[:i]
	}

	i = strings.Index(host, "/")
	if i != -1 {
		host = host[:i]
	}
	requestType = headerMethod[:strings.Index(headerMethod, " ")]

	return requestType, host, protocol, port
}

func (ps *ProxyServer) splitHeader(bytearray []byte) []string {

	result := make([]string, headerLine)
	j := 0
	temp := false

	if(bytearray[0] == 0) {
		fmt.Println("ps: splitHeader: Couldn't get httpheader, zero filter")
		result[0] = string("-1")
		return result
	}

	for index, element := range bytearray {
		if(element == 13) {
			if(bytearray[index+1] == 10) {
				temp = true
			}
		}

		if(temp != true) {
			result[j] += string(element)
		}

		if(element == 10) {
			temp = false
			j += 1
		}
	}

	return result;
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
