all:
	CGO_ENABLED=1 go run -race cm.go main.go mailbox.go message.go page.go process.go records.go
r:
	CGO_ENABLED=1 go run -race cm.go main.go mailbox.go message.go page.go process.go records.go -r
w:
	CGO_ENABLED=1 go run -race cm.go main.go mailbox.go message.go page.go process.go records.go -w

build:
	go build cm.go main.go mailbox.go message.go page.go process.go
clean:
	go clean