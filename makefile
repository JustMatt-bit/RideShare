all: build run

build:
	rm -f ./rideshare-go/rideshare
	rm -f ./rideshare-go/rideshare

run:
	go run ./rideshare-go/main.go