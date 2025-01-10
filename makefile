test:
	go test -v ./...

lint:
	golangci-lint run ./...

gosec:
	gosec -quiet ./...

qa: test lint gosec
