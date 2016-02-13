package ToTo

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	bufferLength = 8192
	headerLine   = 30
)

// ProxyServer stores the configuration for ProxyServer and some functions.
type ToTo struct {
	port string
}

// ask port
func (ps *ToTo) askPort() int {
	port := 0

	fmt.Print("Port : ")
	fmt.Scanf("%d", &port)

	return port
}

// initialize socket server
func (ps *ToTo) init(port int) net.Listener {
	ps.port = ":"
	ps.port += strconv.Itoa(port)

	server, err := net.Listen("tcp", ps.port)

	if server == nil {
		panic("init: port listening error : " + err.Error())
	}

	return server
}

// accept client
func (ps *ToTo) acceptClient(server net.Listener) chan net.Conn {
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

// connect host by socket
func (ps *ToTo) connectHost(client net.Conn) {
	HeaderInfo, Datas := ps.getData(client)

	if (Datas[0] == 0) || (HeaderInfo == "-1") {
		return
	}
	requestType, host, _, port := ps.parseHTTPHeaderMethod(HeaderInfo)

	if port == -1 {
		return
	}

	connectionHost, _ := net.Dial("tcp", host+":"+strconv.Itoa(port))

	if requestType == "CONNECT" {
		connectionHost.Write([]byte("HTTP/1.1 200 Connection established\n"))
	} else {
		connectionHost.Write(Datas)
		go func() {
			io.Copy(connectionHost, client)
		}()
		io.Copy(client, connectionHost)
	}
	client.Close()
	connectionHost.Close()
	return
}

// Get first line of http header
func (ps *ToTo) getData(client net.Conn) (string, []byte) {

	buffer := make([]byte, bufferLength)
	client.Read(buffer)

	return ps.splitHeader(buffer)[0], buffer
}

// parse HTTPHeader from packets
func (ps *ToTo) parseHTTPHeaderMethod(headerMethod string) (string, string, string, int) {
	var (
		requestType string
		host        string
		protocol    string
		port        int
	)
	errorCounter := false

	// ex: GET http://google.com/ HTTP/1.1
	temp := headerMethod[strings.Index(headerMethod, " ")+1:]
	protocol = temp[strings.Index(temp, " ")+1:]

	url := temp[:strings.Index(temp, " ")]

	i := strings.Index(url, "://")

	if i == -1 {
		fmt.Println("Uncorrect URL")
		errorCounter = true
	} else {
		host = url[i+3 : len(url)-1]
	}

	if errorCounter == false {
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

	return "", "", "", -1
}

// split header each lines
func (ps *ToTo) splitHeader(bytearray []byte) []string {

	result := make([]string, headerLine)
	j := 0
	temp := false

	if bytearray[0] == 0 {
		fmt.Println("ps: splitHeader: Couldn't get httpheader, zero filter")
		result[0] = string("-1")
		return result
	}

	for index, element := range bytearray {
		if element == '\r' {
			if bytearray[index+1] == '\n' {
				temp = true
			}
		}

		if temp != true {
			result[j] += string(element)
		}

		if element == '\n' {
			temp = false
			j++
		}
	}

	return result
}
