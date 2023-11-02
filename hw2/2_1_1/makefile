all:
	CGO_ENABLED=1 go run -race main.go message.go node.go queue.go set.go vectorclock.go


build:
	go build main.go message.go node.go queue.go set.go vectorclock.go
clean:
	go clean