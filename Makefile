
local: FORCE
	go build -o ~/.bin/tvecty cmd/tvecty/main.go

test: FORCE
	go test ./...

testr: FORCE
	go test -race ./...

fmt: FORCE
	go fmt ./...

FORCE: