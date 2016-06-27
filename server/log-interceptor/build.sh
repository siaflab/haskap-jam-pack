#!/bin/sh

#BUILD_VERSION=`git describe --always --dirty`
BUILD_VERSION=0.3.0
echo BUILD_VERSION: $BUILD_VERSION
BUILD_DATE=`date +%FT%T%z`
echo BUILD_DATE: $BUILD_DATE
GIT_COMMIT=`git describe --always`
echo GIT_COMMIT: $GIT_COMMIT

SRC_DIR=.
OUT_DIR=bin

mkdir -p $OUT_DIR
export GOOS=darwin;export GOARCH=amd64;go get github.com/google/gopacket;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/haskap-jam-interceptor ${SRC_DIR}/haskap-jam-interceptor.go

# cannot cross-compile pcap...
#  https://stackoverflow.com/questions/31648793/go-programming-cross-compile-for-revel-framework
# export GOOS=windows;export GOARCH=amd64;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/haskap-jam-interceptor ${SRC_DIR}/haskap-jam-interceptor.go
# export GOOS=windows;export GOARCH=386;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/haskap-jam-interceptor ${SRC_DIR}/haskap-jam-interceptor.go
# export GOOS=linux;export GOARCH=amd64;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/haskap-jam-interceptor ${SRC_DIR}/haskap-jam-interceptor.go
# export GOOS=linux;export GOARCH=386;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/haskap-jam-interceptor ${SRC_DIR}/haskap-jam-interceptor.go
# export GOOS=linux;export GOARCH=arm;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/haskap-jam-interceptor ${SRC_DIR}/haskap-jam-interceptor.go
