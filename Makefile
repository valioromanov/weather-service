TERRAFORM_DIR := terraform

.PHONY: build tests testsWithCoverage deploy

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap ./cmd/lambda
	zip lambda.zip bootstrap

tests:
	 ginkgo run ./...

testsWithCoverage:
	ginkgo -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

deploy:
	terraform -chdir=$(TERRAFORM_DIR) init
	terraform -chdir=$(TERRAFORM_DIR) apply -auto-approve