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

	var sem = make(chan int, 1)
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

		go func(b []byte) {
			printTestServerReceivedMessage(rcvFromAddr, string(b))
			sem <- 1
			server.ReceiveMsgList.PushBack(b)
			<-sem
		}(buf[:n])
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
func setup() (*TestServer, *PassthroughServer, *net.UDPConn) {
	testServer := &TestServer{4557, nil, nil}
	go testServer.Start()
	time.Sleep(50 * time.Millisecond)

	rcvPort, sndPort := 4559, 4557
	server := &PassthroughServer{rcvPort, sndPort, nil, nil}
	go server.Start()
	time.Sleep(50 * time.Millisecond)

	testSndConn := prepareSendConnection(4559)
	return testServer, server, testSndConn
}

func teardown(testServer *TestServer, server *PassthroughServer, testSndConn *net.UDPConn) {
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
	time.Sleep(50 * time.Millisecond)
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

	time.Sleep(50 * time.Millisecond)
	sendBytes2 := []byte("val2")
	testSndConn.Write(sendBytes2)

	//// assert
	time.Sleep(50 * time.Millisecond)
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
	sendBytes2 := []byte("val2")
	go testSndConn.Write(sendBytes)
	go testSndConn.Write(sendBytes2)

	//// assert
	time.Sleep(50 * time.Millisecond)
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

func TestSend200MessageSimul(t *testing.T) {
	//// setup
	testServer, server, testSndConn := setup()

	//// test
	const msgNum = 200
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
			testSndConn.Write(sendBytes)
			sendMsgList.PushBack(sendBytes)
		}(i)
	}

	//// assert
	time.Sleep(50 * time.Millisecond)
	listLen := testServer.ReceiveMsgList.Len()
	if listLen != msgNum {
		t.Errorf("%q - Expected %q, actual %q", "testServer.ReceiveMsgList.Len()", 1, listLen)
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
	time.Sleep(50 * time.Millisecond)
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
