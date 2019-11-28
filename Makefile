.PHONY: build clean

build:
	@go version
	GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod=vendor -v \
		-ldflags "-s -w -X main.Version=v$(shell date +%y.%-m.%-d)" \
		-o build/tgmid main.go
	upx -9 build/tgmid
clean:
	rm -rf build
