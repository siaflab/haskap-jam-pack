package main

import (
	"bytes"
	"container/list"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"testing"
	"time"
)

var c chan int

// TestServer represents a server process.
type TestServer struct {
	ReceivePort    int
	ReceiveMsgList *list.List
	rcvConn        *net.UDPConn
}

// Start starts listening.
func (server *TestServer) Start() {
	server.rcvConn = prepareReceiveConnection(server.ReceivePort)
	defer server.rcvConn.Close()

	server.ReceiveMsgList = list.New()

	printTestServerStartedMessage(server.ReceivePort)

	c = make(chan int)
	receiveCount := 0
	for {
		if server.rcvConn == nil {
			break
		}
		buf := make([]byte, 1024*10)
		n, rcvFromAddr, err := server.rcvConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		b := buf[:n]
		printTestServerReceivedMessage(rcvFromAddr, string(b))
		server.ReceiveMsgList.PushBack(b)
		receiveCount++
		c <- receiveCount
	}
}

// Stop stops listening.
func (server *TestServer) Stop() {
	defer func() {
		server.rcvConn.Close()
		server.rcvConn = nil
		time.Sleep(50 * time.Millisecond)
	}()
}

// ReceivedMsgListContainsBytes checks if ReceivedMsgList contains the specified b.
func (server *TestServer) ReceivedMsgListContainsBytes(b []byte) bool {
	l := server.ReceiveMsgList
	for e := l.Front(); e != nil; e = e.Next() {
		v, ok := e.Value.([]byte)
		if !ok {
			fmt.Println("can not convert to []byte:", e)
			continue
		}

		if bytes.Equal(v, b) {
			return true
		}
	}
	return false
}

func prepareReceiveConnection(port int) *net.UDPConn {
	rcvAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port)) // from any address at specified port
	handleError(err)
	rcvConn, err := net.ListenUDP("udp", rcvAddr)
	handleError(err)
	return rcvConn
}

func printTestServerStartedMessage(rcvPort int) {
	fmt.Println("*****")
	fmt.Println(time.Now())
	fmt.Println("TestServer started successfully.")
	fmt.Println("listening to udp", rcvPort)
	fmt.Println("*****")
}

func printTestServerReceivedMessage(rcvFromAddr *net.UDPAddr, message string) {
	fmt.Println("+++++")
	fmt.Println(time.Now())
	fmt.Println("Received from", rcvFromAddr)
	fmt.Println(message)
	fmt.Println("+++++")
}

//// Test utils
func setup() (*TestServer, *Interceptor, *net.UDPConn) {
	testServer := &TestServer{3333, nil, nil}
	go testServer.Start()
	time.Sleep(50 * time.Millisecond)

	// WORKAROUND for travis
	var deviceName string
	if runtime.GOOS == "darwin" {
		deviceName = "lo0" // osx
	} else {
		deviceName = "lo" // travis
	}
	rcvPort, sndPort := 4558, 3333
	sndAddress := "127.0.0.1"
	server := &Interceptor{deviceName, rcvPort, sndAddress, sndPort, nil}
	go server.Start()
	time.Sleep(50 * time.Millisecond)

	testSndConn := prepareSendConnection("127.0.0.1", 4558)
	return testServer, server, testSndConn
}

func teardown(testServer *TestServer, server *Interceptor, testSndConn *net.UDPConn) {
	defer testSndConn.Close()
	defer server.Stop()
	defer testServer.Stop()
	time.Sleep(50 * time.Millisecond)
}

//// Test Cases
func TestStart(t *testing.T) {
	//// setup
	testServer, server, testSndConn := setup()

	//// teardown
	teardown(testServer, server, testSndConn)
}

func TestSend1Message(t *testing.T) {
	//// setup
	testServer, server, testSndConn := setup()

	//// test
	sendBytes := []byte("val")
	testSndConn.Write(sendBytes)

	//// assert
	for {
		receiveCount := <-c
		if receiveCount >= 1 {
			break
		}
	}
	listLen := testServer.ReceiveMsgList.Len()
	if listLen != 1 {
		t.Errorf("%q - Expected %q, actual %q", "testServer.ReceiveMsgList.Len()", 1, listLen)
	}

	first := testServer.ReceiveMsgList.Front()
	expected1, _ := first.Value.([]byte)
	if !bytes.Equal(expected1, sendBytes) {
		t.Errorf("Expected %q, actual %q", sendBytes, expected1)
	}

	//// teardown
	teardown(testServer, server, testSndConn)
}
