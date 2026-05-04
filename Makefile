BINARY ?= ocm-csv-parser

.PHONY: build test fmt vet clean

build:
	go build -o $(BINARY) ./...

test:
	go test ./... -v

fmt:
	gofmt -w .

vet:
	go vet ./...

clean:
	rm -f $(BINARY)
