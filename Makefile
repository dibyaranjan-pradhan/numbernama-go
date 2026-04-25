APP=numbernama-go
CMD=./cmd

.PHONY: all build run tidy wire test docker-build clean

all: build

build:
	go build -o bin/$(APP) $(CMD)

run: build
	PORT=7002 ./bin/$(APP)

tidy:
	go mod tidy

wire:
	go run github.com/google/wire/cmd/wire@latest $(CMD)

test:
	go test ./...

docker-build:
	docker build -t $(APP):local .

clean:
	rm -rf bin/
