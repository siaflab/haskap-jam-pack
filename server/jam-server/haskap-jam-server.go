//
// Copyright (c) 2016 SIAF LAB.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	maxOutstanding     = 1
	bufferLen          = 1024 * 10
	defaultReceivePort = 4559
	defaultSendPort    = 4557
)

var (
	// BuildVersion sets version string
	BuildVersion string

	// GitCommit sets commit hash of git
	GitCommit string

	// BuildDate sets date of built datetime
	BuildDate string

	debugMode bool = true
)

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}

// Config
var config struct {
	ReceivePort int `json:"receivePort"`
	SocicPiPort int `json:"socicPiPort"`
}

func readConfig() (int, int) {
	rcvPort := defaultReceivePort
	sndPort := defaultSendPort
	configFile, err := os.Open("haskap-jam-server-config.json")
	if err != nil {
		fmt.Println("Error - opening config file:", err.Error())
		return rcvPort, sndPort
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		fmt.Println("Error - parsing config file:", err.Error())
		return rcvPort, sndPort
	}

	fmt.Println("config.ReceivePort:", config.ReceivePort)
	fmt.Println("config.SocicPiPort:", config.SocicPiPort)

	if config.ReceivePort != 0 {
		rcvPort = config.ReceivePort
	}
	if config.SocicPiPort != 0 {
		sndPort = config.SocicPiPort
	}

	return rcvPort, sndPort
}

// PassthroughServer represents a server process.
type PassthroughServer struct {
	ReceivePort int
	SendPort    int
	rcvListener *net.TCPListener
	sndConn     *net.UDPConn
}

// Start starts listening.
func (server *PassthroughServer) Start() {
	server.sndConn = prepareSendConnection(server.SendPort)
	defer server.sndConn.Close()

	rcvAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(server.ReceivePort)) // from any address at specified port
	handleError(err)
	server.rcvListener, err = net.ListenTCP("tcp", rcvAddr)
	handleError(err)

	printStartedMessage(server.ReceivePort, server.SendPort)

	var sem = make(chan int, maxOutstanding)
	for {
		if server.rcvListener == nil || server.sndConn == nil {
			break
		}

		conn, err := server.rcvListener.Accept()
		if err != nil || conn == nil {
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			err := conn.SetDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			buf := make([]byte, bufferLen)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			b := buf[:n]
			sem <- 1
			printReceivedMessage(conn.RemoteAddr(), n, string(b))
			n, err = server.sndConn.Write(b)
			if err != nil {
				fmt.Println("Error:", err)
			}
			<-sem
		}(conn)
	}
}

// Stop stops listening.
func (server *PassthroughServer) Stop() {
	defer func() {
		server.rcvListener.Close()
		server.rcvListener = nil
	}()
	defer func() {
		server.sndConn.Close()
		server.sndConn = nil
	}()
}

func prepareSendConnection(port int) *net.UDPConn {
	sndToAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1"+":"+strconv.Itoa(port))
	sndFromAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	sndConn, err := net.DialUDP("udp", sndFromAddr, sndToAddr)
	handleError(err)
	return sndConn
}

func printStartedMessage(rcvPort int, sndPort int) {
	fmt.Println("#####")
	fmt.Println(time.Now())
	fmt.Println("haskap-jam-server started successfully.")
	fmt.Println("version: " + BuildVersion + ", build: " + GitCommit + ", date:" + BuildDate)
	fmt.Println("listening to udp", rcvPort)
	fmt.Println("and will send to udp", sndPort)
	fmt.Println("#####")
}

func printReceivedMessage(rcvFromAddr net.Addr, n int, message string) {
	// if !debugMode {
	//     return
	// }

	fmt.Println("-----")
	fmt.Println(time.Now())
	fmt.Println("Received from:", rcvFromAddr)
	fmt.Println("size:", n)
	fmt.Println(message)
	fmt.Println("-----")
}

// main
func printVersion() {
	fmt.Println("haskap-jam-server version " + BuildVersion + ", build " + GitCommit + ", date " + BuildDate)
}

func main() {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "show version")
	flag.BoolVar(&versionFlag, "v", false, "show version")
	flag.BoolVar(&debugMode, "debug", false, "print debug output")
	flag.BoolVar(&debugMode, "d", false, "print debug output")
	flag.Parse()

	if versionFlag {
		printVersion()
		return
	}

	rcvPort, sndPort := readConfig()
	server := &PassthroughServer{rcvPort, sndPort, nil, nil}
	// server := &PassthroughServer{rcvPort, sndPort, nil}
	server.Start()
}
