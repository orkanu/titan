# Taks to run and build titan

run:
	go run ./cmd/titan.go

run-f:
	go run ./cmd/titan.go fetch

run-i:
	go run ./cmd/titan.go install

run-b:
	go run ./cmd/titan.go build

run-c:
	go run ./cmd/titan.go clean

run-a:
	go run ./cmd/titan.go all

run-s:
	go run ./cmd/titan.go serve -p local:all

run-h:
	go run ./cmd/titan.go help

release:
	./scripts/build.sh

clean:
	rm -rf ./bin
