

BIN_DIR:=/bin
build: 
	env GOOS=darwin GOARCH=amd64 go build -o bin/tfstate-lookup

.PHONY: build