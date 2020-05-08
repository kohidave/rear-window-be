all: build

build:
	go build -o emoji-detective ./pkg
run:
	go run pkg/*
