# Makefile
BINARY_NAME=omgsetup

build:
	go build -o cmd/main/${BINARY_NAME} cmd/main/main.go
