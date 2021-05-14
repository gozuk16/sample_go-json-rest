all:
	go build -ldflags "-X main.version=0.1.0"
run:
	go run -ldflags "-X main.version=0.1.0" main.go
