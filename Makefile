.PHONY:
.SILENT:

build: 
	go build -o ./.bin/marketplace-service ./cmd/app/main.go

run: build 
	./.bin/marketplace-service