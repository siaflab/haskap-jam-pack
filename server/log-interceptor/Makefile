.PHONY: all
all:
	gofmt -w .
	./build.sh

.PHONY: run
run:
	./run.sh

.PHONY: test
test:
	gofmt -w .
	test -z "$(gofmt -s -l . | tee /dev/stderr)"
	test -z "$($GOPATH/bin/golint ./... | tee /dev/stderr)"
	./test.sh

.PHONY: clean
clean:
	rm -rf bin
