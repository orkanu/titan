# Taks to run and build titan

run:
	go run ./cmd/titan.go

run-fetch:
	go run ./cmd/titan.go fetch

release:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o ./build/linux/titan ./cmd/titan.go
	GOOS=linux GOARCH=arm go build -ldflags="-s -w" -trimpath -o ./build/linux/titan-arm ./cmd/titan.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o ./build/windows/titan.exe ./cmd/titan.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o ./build/mac/titan ./cmd/titan.go
