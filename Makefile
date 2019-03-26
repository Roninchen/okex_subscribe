V=okexv11_8
PKG=main.go
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(V) $(PKG)