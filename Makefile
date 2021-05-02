
local: FORCE
	go build -o ~/.bin/tvecty cmd/tvecty/main.go

test: FORCE
	go test ./...

fmt: FORCE
	go fmt ./...

FORCE: