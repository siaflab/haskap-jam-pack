#!/bin/sh

#BUILD_VERSION=`git describe --always --dirty`
version=`cat ../../../version.txt`
BUILD_VERSION=$version
echo BUILD_VERSION: $BUILD_VERSION
BUILD_DATE=`date +%FT%T%z`
echo BUILD_DATE: $BUILD_DATE
GIT_COMMIT=`git describe --always`
echo GIT_COMMIT: $GIT_COMMIT

SRC_DIR=.
OUT_DIR=bin
SRC_FILE=test-client.go
OUT_FILE=test-client

mkdir -p $OUT_DIR

export GOOS=darwin;export GOARCH=amd64;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/${OUT_FILE} ${SRC_DIR}/${SRC_FILE}
# export GOOS=windows;export GOARCH=amd64;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/${OUT_FILE} ${SRC_DIR}/${SRC_FILE}
# export GOOS=windows;export GOARCH=386;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/${OUT_FILE} ${SRC_DIR}/${SRC_FILE}
# export GOOS=linux;export GOARCH=amd64;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/${OUT_FILE} ${SRC_DIR}/${SRC_FILE}
# export GOOS=linux;export GOARCH=386;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/${OUT_FILE} ${SRC_DIR}/${SRC_FILE}
# export GOOS=linux;export GOARCH=arm;go build -ldflags "-X main.BuildVersion=$BUILD_VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildDate=$BUILD_DATE" -o ${OUT_DIR}/${GOOS}_${GOARCH}/${OUT_FILE} ${SRC_DIR}/${SRC_FILE}
