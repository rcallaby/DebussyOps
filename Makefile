.PHONY: build run fmt test

build:
	go build -o bin/orchestrator ./cmd/orchestrator
	go build -o bin/calendar ./examples/agents/calendar
	go build -o bin/todo ./examples/agents/todo

run: build
	./bin/orchestrator &
	sleep 1; ./bin/calendar &
	sleep 1; ./bin/todo &

fmt:
	gofmt -w .

test:
	go test ./...