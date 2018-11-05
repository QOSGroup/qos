########################################
### Build/Install

build:
ifeq ($(OS),Windows_NT)
	go build -o build/qosd.exe ./cmd/qosd
	go build -o build/qoscli.exe ./cmd/qoscli
else
	go build -o build/qosd ./cmd/qosd
	go build -o build/qoscli ./cmd/qoscli
endif

install:
	go install ./cmd/qosd
	go install ./cmd/qoscli