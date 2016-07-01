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

// osc encoding process is taken from https://github.com/aike/oscer/blob/master/src/osc/osc.go

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
)

const (
	defaultServerIp   = "127.0.0.1"
	defaultServerPort = "4559"
	defaultCode       = `#load "~/github/haskap-jam-pack/client/haskap-jam-loop.rb"

live_loop :test do
play 60
sleep 1
end`
)

var (
	// BuildVersion sets version string
	BuildVersion string
	// GitCommit sets commit hash of git
	GitCommit string
	// BuildDate sets date of built datetime
	BuildDate string
	senddata  []byte
	oscarg    []byte
	initdata  bool = false
)

func send(serverIP, serverPort string) {
	if !initdata {
		return
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", serverIP+":"+serverPort)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer conn.Close()

	conn.Write(senddata)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
		os.Exit(1)
	}
}

func pushOscArgs(arr []string) error {
	for i := 0; i < len(arr); i++ {
		if match(`^[+-]?[0-9]+$`, arr[i]) {
			// Int32
			num_i64, err := strconv.ParseInt(arr[i], 10, 32)
			num_i32 := int32(num_i64)
			if err != nil {
				return errors.New("osc args error")
			}
			pushDataI32(num_i32)

		} else if match(`^[+-]?[0-9.]+$`, arr[i]) {
			// Float32
			num_f64, err := strconv.ParseFloat(arr[i], 32)
			num_f32 := float32(num_f64)
			if err != nil {
				return errors.New("osc args error")
			}
			pushDataF32(num_f32)

		} else {
			// String
			pushDataString(arr[i])
		}
	}

	senddata = append(senddata, 0)
	fill4byte()
	senddata = append(senddata, oscarg...)

	initdata = true
	return nil
}

func match(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}

func pushOscAddress(str string) {
	senddata = append(senddata, []byte(str)...)
	senddata = append(senddata, 0)
	fill4byte()
	senddata = append(senddata, 0x2c)
}

func fill4byte() {
	for datalen := len(senddata); datalen%4 != 0; datalen++ {
		senddata = append(senddata, 0)
	}
}

func pushDataI32(num int32) {
	senddata = append(senddata, 'i')
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, num)
	oscarg = append(oscarg, buf.Bytes()...)
}

func pushDataF32(num float32) {
	senddata = append(senddata, 'f')
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, num)
	oscarg = append(oscarg, buf.Bytes()...)
}

func pushDataString(str string) {
	senddata = append(senddata, 's')
	buf := bytes.NewBuffer([]byte(str))
	oscarg = append(oscarg, buf.Bytes()...)

	oscarg = append(oscarg, 0)
	for datalen := len(oscarg); datalen%4 != 0; datalen++ {
		oscarg = append(oscarg, 0)
	}
}

func main() {
	// $ go run test-client.go 192.168.100.139
	// $ go run test-client.go 192.168.100.139 4559
	// $ go run test-client.go 192.168.100.139 4559 "play 60"
	fmt.Println("version: " + BuildVersion + ", build: " + GitCommit + ", date:" + BuildDate)

	argslen := len(os.Args)
	if argslen < 1 {
		fmt.Fprintf(os.Stderr, "usage: test-client [serverIP] [serverPort] [code]\n")
		os.Exit(1)
		return
	}

	serverIP := defaultServerIp
	if argslen > 1 && len(os.Args[1]) != 0 {
		serverIP = os.Args[1]
	}

	serverPort := defaultServerPort
	if argslen > 2 && len(os.Args[2]) != 0 {
		serverPort = os.Args[2]
	}

	senddata = []byte{}
	oscarg = []byte{}

	pushOscAddress("/save-and-run-buffer")

	clientIP := "127.0.0.1"
	bufferNo := "one"

	clientId := "haskap-test-client-" + clientIP
	bufferId := "haskap-test-buffer-" + clientIP + "-" + bufferNo
	code := defaultCode
	if argslen > 3 && len(os.Args[3]) != 0 {
		code = os.Args[3]
	}
	workspaceId := "haskap-test-workspace-" + clientIP + "-" + bufferNo

	arr := []string{
		clientId,
		bufferId,
		code,
		workspaceId}

	pushOscArgs(arr)

	for i := 0; i < len(senddata); i++ {
		fmt.Printf("%02x ", senddata[i])
		if i%16 == 15 {
			fmt.Printf("\n")
		}
	}

	send(serverIP, serverPort)
}
