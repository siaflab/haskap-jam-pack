#!/bin/sh

echo $GOPATH
go_path=$GOPATH

go get github.com/google/gopacket

go_cmd=`which go`
echo $go_cmd
sudo GOPATH=$go_path $go_cmd test
