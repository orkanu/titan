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

release:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o ./build/linux/titan ./cmd/titan.go
	GOOS=linux GOARCH=arm go build -ldflags="-s -w" -trimpath -o ./build/linux/titan-arm ./cmd/titan.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o ./build/windows/titan.exe ./cmd/titan.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o ./build/mac/titan ./cmd/titan.go
