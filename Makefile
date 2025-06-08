TARGET := build/sem

MAIN := ./cmd/sem

build:
	go build -o $(TARGET) $(MAIN)

build-win:
	set GOOS=windows
	set GOARCH=amd64
	set CGO_ENABLED=0
	go build -o $(TARGET).exe $(MAIN)
