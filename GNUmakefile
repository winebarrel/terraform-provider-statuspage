.PHONY: build vet test lint clean docs

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

docs:
	go generate ./...

dev.tfrc:
	sed "s|{{PATH_TO_PROVIDER}}|$(shell pwd)|" dev.tfrc.tpl > dev.tfrc

.PHONY: tf-plan
tf-plan: build dev.tfrc
	TF_CLI_CONFIG_FILE=dev.tfrc terraform plan

.PHONY: tf-apply
tf-apply: build dev.tfrc
	TF_CLI_CONFIG_FILE=dev.tfrc terraform apply -auto-approve

.PHONY: tf-clean
tf-clean: clean
	rm -f dev.tfrc terraform.tfstate*
