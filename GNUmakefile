.PHONY: build vet test testacc lint clean

build: vet
	go build

vet:
	go vet ./...

test:
	go test -v -count=1 ./...

testacc:
	TF_ACC=1 go test -v -count=1 -timeout 120m ./internal/provider/

lint:
	golangci-lint run

clean:
	rm -f terraform-provider-statuspage
