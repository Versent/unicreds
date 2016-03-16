NAME=unicreds
ARCH=$(shell uname -m)
VERSION=1.0.5
GO15VENDOREXPERIMENT := 1
ITERATION := 1

build:
	rm -rf build && mkdir build
	mkdir -p build/Linux  && GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o build/Linux/$(NAME) ./cmd/unicreds
	mkdir -p build/Darwin && GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o build/Darwin/$(NAME) ./cmd/unicreds
	mkdir -p build/Darwin && GOOS=windows go build -ldflags "-X main.Version=$(VERSION)" -o build/Windows/$(NAME).exe ./cmd/unicreds


fmt:
	gofmt -w=true $$(find . -type f -name '*.go')
	goimports -w=true -d $$(find . -type f -name '*.go')

test:
	go test -v ./...

release: build
	git push origin master
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_$(VERSION)_linux_$(ARCH).tgz -C build/Linux $(NAME)
	tar -zcf release/$(NAME)_$(VERSION)_darwin_$(ARCH).tgz -C build/Darwin $(NAME)
	tar -zcf release/$(NAME)_$(VERSION)_windows_$(ARCH).tgz -C build/Windows $(NAME).exe
	gh-release create versent/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD)

packages:
	rm -rf package && mkdir package
	rm -rf stage && mkdir -p stage/usr/bin
	cp build/Linux/unicreds stage/usr/bin
	fpm --name $(NAME) -a x86_64 -t rpm -s dir --version $(VERSION) --iteration $(ITERATION) -C stage -p package/$(NAME)-$(VERSION)_$(ITERATION).rpm usr

.PHONY: build fmt test release packages
