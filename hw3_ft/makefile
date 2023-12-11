all:
	CGO_ENABLED=1 go run -race cm.go cm_ref.go main.go mailbox.go message.go records.go page.go process.go role.go
r:
	CGO_ENABLED=1 go run -race cm.go cm_ref.go main.go mailbox.go message.go records.go page.go process.go role.go -r=0
w:
	CGO_ENABLED=1 go run -race cm.go cm_ref.go main.go mailbox.go message.go records.go page.go process.go role.go -w=100

build:
	go build cm.go cm_ref.go main.go mailbox.go message.go records.go page.go process.go role.go
clean:
	go clean