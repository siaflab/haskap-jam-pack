.PHONY: all
all:
	gofmt -w .
	./build.sh

.PHONY: run
run:
	go run spi-ping.go

# .PHONY: test
# test:
# 	gofmt -w .
# 	test -z "$(gofmt -s -l . | tee /dev/stderr)"
# 	test -z "$($GOPATH/bin/golint ./... | tee /dev/stderr)"
# 	go test

.PHONY: clean
clean:
	rm -rf bin
