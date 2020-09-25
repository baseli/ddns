APP=ddns
.PHONY: build
build:
	GOOS=darwin GOARCH=amd64 go build -o build/${APP}_darwin_amd64 cmd/main.go
	GOOS=linux GOARCH=amd64 go build -o build/${APP}_linux_amd64 cmd/main.go
	GOOS=linux GOARCH=arm go build -o build/${APP}_linux_arm cmd/main.go