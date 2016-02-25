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
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	defaultDeviceName    = "lo0"
	defaultReceivePort   = 4558
	defaultSendToAddress = "127.0.0.1"
	defaultSendToPort    = 3333
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
	DeviceName    string `json:"deviceName"`
	ReceivePort   int    `json:"receivePort"`
	SendToAddress string `json:"sendToAddress"`
	SendToPort    int    `json:"sendToPort"`
}

func readConfig() (string, int, string, int) {
	deviceName := defaultDeviceName
	rcvPort := defaultReceivePort
	sndAddress := defaultSendToAddress
	sndPort := defaultSendToPort
	configFile, err := os.Open("haskap-jam-interceptor-config.json")
	if err != nil {
		fmt.Println("Error - opening config file:", err.Error())
		return deviceName, rcvPort, sndAddress, sndPort
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		fmt.Println("Error - parsing config file:", err.Error())
		return deviceName, rcvPort, sndAddress, sndPort
	}

	fmt.Println("config.DeviceName:", config.DeviceName)
	fmt.Println("config.ReceivePort:", config.ReceivePort)
	fmt.Println("config.SendToAddress:", config.SendToAddress)
	fmt.Println("config.SendToPort:", config.SendToPort)

	if len(config.DeviceName) > 0 {
		deviceName = config.DeviceName
	}
	if config.ReceivePort != 0 {
		rcvPort = config.ReceivePort
	}
	if len(config.SendToAddress) > 0 {
		sndAddress = config.SendToAddress
	}
	if config.SendToPort != 0 {
		sndPort = config.SendToPort
	}

	return deviceName, rcvPort, sndAddress, sndPort
}

// Interceptor represents a server process.
type Interceptor struct {
	ReceiveDeviceName string
	ReceivePort       int
	SendAddress       string
	SendPort          int
	sndConn           *net.UDPConn
}

// Start starts listening.
func (server *Interceptor) Start() {
	const (
		snapshotLen int32         = 1024
		promiscuous bool          = false
		timeout     time.Duration = 10 * time.Millisecond
	)

	// sender
	server.sndConn = prepareSendConnection(server.SendAddress, server.SendPort)

	// Open device
	device := server.ReceiveDeviceName
	handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Set filter
	watchPort := server.ReceivePort
	var filter string = fmt.Sprintf("udp and port %d", watchPort)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	printStartedMessage(server.ReceivePort, server.SendAddress, server.SendPort)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Do something with a packet here.
		processPacket(packet, server.sndConn)
	}
}

// Stop stops listening.
func (server *Interceptor) Stop() {
	defer func() {
		server.sndConn.Close()
		server.sndConn = nil
	}()
}

func prepareSendConnection(addr string, port int) *net.UDPConn {
	sndToAddr, err := net.ResolveUDPAddr("udp", addr+":"+strconv.Itoa(port))
	sndFromAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	sndConn, err := net.DialUDP("udp", sndFromAddr, sndToAddr)
	handleError(err)
	return sndConn
}

func printStartedMessage(watchPort int, sndToAddr string, sndToPort int) {
	fmt.Println("#####")
	fmt.Println(time.Now())
	fmt.Println("haskap-jam-interceptor started successfully.")
	fmt.Println("version: " + BuildVersion + ", build: " + GitCommit + ", date:" + BuildDate)
	fmt.Printf("capturing UDP port %d packets.", watchPort)
	fmt.Println()
	fmt.Printf("and will send to %s:%d", sndToAddr, sndToPort)
	fmt.Println()
	fmt.Println("#####")
}

func processPacket(packet gopacket.Packet, sndConn *net.UDPConn) {
	if debugMode {
		fmt.Println("-----")
		// Let's see if the packet is IP (even though the ether type told us)
		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
			fmt.Println("IPv4 layer detected.")
			ip, _ := ipLayer.(*layers.IPv4)

			// IP layer variables:
			// Version (Either 4 or 6)
			// IHL (IP Header Length in 32-bit words)
			// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
			// Checksum, SrcIP, DstIP
			fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
			fmt.Println("Protocol: ", ip.Protocol)
			fmt.Println()
		}

		// Iterate over all layers, printing out each layer type
		fmt.Println("All packet layers:")
		for _, layer := range packet.Layers() {
			fmt.Println("- ", layer.LayerType())
		}
	}

	// Check for errors
	if err := packet.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}

	// Let's see if the packet is UDP
	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer == nil {
		// not UDP Protocol
		return
	}

	if debugMode {
		fmt.Println("UDP layer detected.")
	}

	// When iterating through packet.Layers() above,
	// if it lists Payload layer then that is the same as
	// this applicationLayer. applicationLayer contains the payload
	if applicationLayer := packet.ApplicationLayer(); applicationLayer != nil {
		buffer := applicationLayer.Payload()
		if debugMode {
			fmt.Println("Application layer/Payload found.")
			fmt.Printf("%s\n", buffer)
		}
		if sndConn != nil {
			_, err := sndConn.Write(buffer)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
	}
}

// main
func printVersion() {
	fmt.Println("haskap-jam-interceptor version " + BuildVersion + ", build " + GitCommit + ", date " + BuildDate)
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

	// deviceName := "lo0"
	// rcvPort, sndPort := 4558, 4560
	// sndAddress := "127.0.0.1"
	deviceName, rcvPort, sndAddress, sndPort := readConfig()
	server := &Interceptor{deviceName, rcvPort, sndAddress, sndPort, nil}
	server.Start()
}
