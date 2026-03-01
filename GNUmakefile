.PHONY: build vet test lint clean

build: vet
	go build

vet:
	go vet ./...

test:
	go test -v -count=1 ./...

lint:
	golangci-lint run

clean:
	rm -f terraform-provider-statuspage
