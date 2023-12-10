all:
	CGO_ENABLED=1 go run -race cm.go cm_ref.go main.go mailbox.go message.go mode.go page.go process.go role.go
r:
	CGO_ENABLED=1 go run -race cm.go cm_ref.go main.go mailbox.go message.go mode.go page.go process.go role.go -r
w:
	CGO_ENABLED=1 go run -race cm.go cm_ref.go main.go mailbox.go message.go mode.go page.go process.go role.go -w

build:
	go build cm.go main.go mailbox.go message.go page.go process.go
clean:
	go clean