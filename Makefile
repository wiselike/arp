all:
	GOARCH=mipsle go build -o arp -ldflags "-s -w" *.go

test:
	go test -v ./...
