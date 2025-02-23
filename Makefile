PROJECT_NAME=kryptonim-example

VERSION=$(shell git describe --always --tags)
CONTAINER_NAME=tomcyr/$(PROJECT_NAME)
IMAGE_TAG=$(CONTAINER_NAME):$(VERSION)
CLI_PATH="cmd/api/main.go"

build:
	go build -o api $(CLI_PATH)

build-image:
	docker build . -t $(IMAGE_TAG)

run-image:
	docker run --env-file .env -p 8080:8080 $(IMAGE_TAG)

test:
	go test ./...