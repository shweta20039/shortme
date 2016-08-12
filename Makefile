dep:
	go get ./...

test:
	go test -v ./...

fmt: 
	find . -type f -name "*.go" | xargs gofmt -s -w

build: dep test fmt
	go build -o shortme main.go

clean:
	rm -f shortme

.PHONY: fmt test dep build clean
