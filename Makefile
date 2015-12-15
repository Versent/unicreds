NAME=unicreds
ARCH=$(shell uname -m)
VERSION=1.0.2
GO15VENDOREXPERIMENT := 1

vendor:
	godep save -d -t

build:
	rm -rf build && mkdir build
	mkdir -p build/Linux  && GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o build/Linux/$(NAME) ./cmd/unicreds
	mkdir -p build/Darwin && GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o build/Darwin/$(NAME) ./cmd/unicreds
	mkdir -p build/Darwin && GOOS=windows go build -ldflags "-X main.Version=$(VERSION)" -o build/Windows/$(NAME).exe ./cmd/unicreds

test:
	go test ./...

release: build
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_$(VERSION)_linux_$(ARCH).tgz -C build/Linux $(NAME)
	tar -zcf release/$(NAME)_$(VERSION)_darwin_$(ARCH).tgz -C build/Darwin $(NAME)
	tar -zcf release/$(NAME)_$(VERSION)_windows_$(ARCH).tgz -C build/Windows $(NAME).exe
	gh-release create versent/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD)

.PHONY: vendor build test release
