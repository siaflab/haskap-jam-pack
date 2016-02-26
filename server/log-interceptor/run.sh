#!/bin/sh

echo $GOPATH
go_path=$GOPATH

go get github.com/google/gopacket
sudo GOPATH=$go_path go run haskap-jam-interceptor.go
