BUILD_DIR := build
TARGET := sem
MAIN := ./cmd/sem

build: build-winows-amd64 build-linux-amd64 build-linux-arm64 build-macos-intel64 build-macos-arm64

build-winows-amd64:
	set GOOS=windows
	set GOARCH=amd64
	set CGO_ENABLED=0
	go build -o ./$(BUILD_DIR)/windows-amd64/$(TARGET).exe $(MAIN)

build-linux-amd64:
	set GOOS=linux
	set GOARCH=amd64
	set CGO_ENABLED=0
	go build -o ./$(BUILD_DIR)/linux-amd64/$(TARGET) $(MAIN)

build-linux-arm64:
	set GOOS=linux
	set GOARCH=arm64
	set CGO_ENABLED=0
	go build -o ./$(BUILD_DIR)/linux-arm64/$(TARGET) $(MAIN)

build-macos-intel64:
	set GOOS=darwin
	set GOARCH=arm64
	set CGO_ENABLED=0
	go build -o ./$(BUILD_DIR)/macos-intel64/$(TARGET) $(MAIN)

build-macos-arm64:
	set GOOS=darwin
	set GOARCH=arm64
	set CGO_ENABLED=0
	go build -o ./$(BUILD_DIR)/macos-arm64/$(TARGET) $(MAIN)
