package main

import (
	"bytes"
	"container/list"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"
)

var c chan int

func prepareReceiveUDPConnection(port int) *net.UDPConn {
	rcvAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port)) // from any address at specified port
	handleError(err)
	rcvConn, err := net.ListenUDP("udp", rcvAddr)
	handleError(err)
	return rcvConn
}

func prepareSendTCPConnection(port int) *net.TCPConn {
	sndToAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1"+":"+strconv.Itoa(port))
	sndFromAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	sndConn, err := net.DialTCP("tcp", sndFromAddr, sndToAddr)
	handleError(err)
	return sndConn
}

// TestServer represents a server process.
type TestServer struct {
	ReceivePort    int
	ReceiveMsgList *list.List
	rcvConn        *net.UDPConn
}

// Start starts listening.
func (server *TestServer) Start() {
	server.rcvConn = prepareReceiveUDPConnection(server.ReceivePort)
	defer server.rcvConn.Close()

	server.ReceiveMsgList = list.New()

	printTestServerStartedMessage(server.ReceivePort)

	c = make(chan int)
	recievCount := 0
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

		go func(buffer []byte, n int) {
			b := buffer[:n]
			printTestServerReceivedMessage(rcvFromAddr, string(b))
			server.ReceiveMsgList.PushBack(b)
			recievCount++
			c <- recievCount
		}(buf, n)
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
func setup() (*TestServer, *PassthroughServer, *net.TCPConn) {
	testServer := &TestServer{4557, nil, nil}
	go testServer.Start()
	time.Sleep(50 * time.Millisecond)

	rcvPort, sndPort := 4559, 4557
	server := &PassthroughServer{rcvPort, sndPort, nil, nil}
	go server.Start()
	time.Sleep(50 * time.Millisecond)

	testSndConn := prepareSendTCPConnection(4559)
	return testServer, server, testSndConn
}

func teardown(testServer *TestServer, server *PassthroughServer, testSndConn *net.TCPConn) {
	defer testSndConn.Close()
	defer server.Stop()
	defer testServer.Stop()
	time.Sleep(50 * time.Millisecond)
}

// Test Cases
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

func TestSend2Message(t *testing.T) {
	//// setup
	testServer, server, testSndConn := setup()

	//// test
	sendBytes := []byte("val")
	testSndConn.Write(sendBytes)
	testSndConn.Close()

	time.Sleep(50 * time.Millisecond)
	sendBytes2 := []byte("val2")
	testSndConn = prepareSendTCPConnection(4559)
	testSndConn.Write(sendBytes2)

	//// assert
	for {
		receiveCount := <-c
		if receiveCount >= 2 {
			break
		}
	}
	listLen := testServer.ReceiveMsgList.Len()
	if listLen != 2 {
		t.Errorf("%q - Expected %q, actual %q", "testServer.ReceiveMsgList.Len()", 2, listLen)
	}

	first := testServer.ReceiveMsgList.Front()
	expected1, _ := first.Value.([]byte)
	if !bytes.Equal(expected1, sendBytes) {
		t.Errorf("Expected %q, actual %q", sendBytes, expected1)
	}

	second := first.Next()
	expected2, _ := second.Value.([]byte)
	if !bytes.Equal(expected2, sendBytes2) {
		t.Errorf("Expected %q, actual %q", sendBytes2, expected2)
	}

	//// teardown
	teardown(testServer, server, testSndConn)
}

func TestSend2MessageSimul(t *testing.T) {
	//// setup
	testServer, server, testSndConn := setup()

	//// test
	sendBytes := []byte("val")
	go func(sendBytes []byte) {
		sndConn := prepareSendTCPConnection(4559)
		sndConn.Write(sendBytes)
		sndConn.Close()
	}(sendBytes)
	sendBytes2 := []byte("val2")
	go func(sendBytes []byte) {
		sndConn := prepareSendTCPConnection(4559)
		sndConn.Write(sendBytes)
		sndConn.Close()
	}(sendBytes2)

	//// assert
	for {
		receiveCount := <-c
		if receiveCount >= 2 {
			break
		}
	}
	listLen := testServer.ReceiveMsgList.Len()
	if listLen != 2 {
		t.Errorf("%q - Expected %q, actual %q", "testServer.ReceiveMsgList.Len()", 2, listLen)
	}

	if !testServer.ReceivedMsgListContainsBytes(sendBytes) {
		t.Errorf("sent bytes not foud: %q", sendBytes)
	}
	if !testServer.ReceivedMsgListContainsBytes(sendBytes2) {
		t.Errorf("sent bytes not foud: %q", sendBytes2)
	}

	//// teardown
	teardown(testServer, server, testSndConn)
}

func TestSend10MessageSimul(t *testing.T) {
	//// setup
	testServer, server, testSndConn := setup()

	//// test
	const msgNum = 10
	const codeMsg = `# Welcome to Sonic Pi v2.9

#load "~/github/haskap-jam-pack/client/haskap-jam-loop.rb"

jam_loop :test do
sample :perc_bell, rate: rrand(-1.5, 1.5)
sleep rrand(0.1, 2)
stop
end`
	sendMsgList := list.New()
	for i := 0; i < msgNum; i++ {

		go func(idx int) {
			sendBytes := []byte(codeMsg)
			strconv.AppendInt(sendBytes, int64(idx), 10)
			sndConn := prepareSendTCPConnection(4559)
			sndConn.Write(sendBytes)
			sendMsgList.PushBack(sendBytes)
			fmt.Println("Sent:", idx)
			sndConn.Close()
		}(i)
	}

	//// assert
	for {
		receiveCount := <-c
		fmt.Println("Recept:", receiveCount)
		if receiveCount >= msgNum {
			break
		}
	}
	listLen := testServer.ReceiveMsgList.Len()
	if listLen != msgNum {
		t.Errorf("%q - Expected %d, actual %d", "testServer.ReceiveMsgList.Len()", msgNum, listLen)
	}

	for e := sendMsgList.Front(); e != nil; e = e.Next() {
		sendBytes, ok := e.Value.([]byte)
		if !ok {
			fmt.Println("can not convert to []byte:", e)
			continue
		}

		if !testServer.ReceivedMsgListContainsBytes(sendBytes) {
			t.Errorf("sent bytes not foud: %q", sendBytes)
		}
	}

	//// teardown
	teardown(testServer, server, testSndConn)
}

func TestSend9216byte(t *testing.T) {
	//// setup
	testServer, server, testSndConn := setup()

	//// test
	// bytesLen := 131071  // ng
	// bytesLen := 13107 // ng
	// bytesLen := 8192  // ok
	// bytesLen := 1024 * 10 // ng
	bytesLen := 1024 * 9 // ok
	sendBytes := make([]byte, bytesLen)
	for i, _ := range sendBytes {
		sendBytes[i] = 0x31
	}
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

	if !testServer.ReceivedMsgListContainsBytes(sendBytes) {
		t.Errorf("sent bytes not foud: %q", sendBytes)
	}

	//// teardown
	teardown(testServer, server, testSndConn)
}
