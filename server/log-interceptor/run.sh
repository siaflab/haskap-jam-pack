#!/bin/sh

echo $GOPATH
gpath=$GOPATH

go get github.com/google/gopacket
sudo GOPATH=$gpath go run haskap-jam-interceptor.go
