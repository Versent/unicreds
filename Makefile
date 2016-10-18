NAME=unicreds
ARCH=$(shell uname -m)
VERSION=1.5.1
GO15VENDOREXPERIMENT := 1
ITERATION := 1

test:
	go test -cover -v ./...

integration:
	go test -v integration/integration_test.go

compile:
	@rm -rf build/
	@gox -ldflags "-X main.Version=$(VERSION)" \
	-osarch="darwin/amd64" \
	-osarch="linux/amd64" \
	-osarch="windows/amd64" \
	-output "build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)" \
	./...

dist: compile
	$(eval FILES := $(shell ls build))
	@rm -rf dist && mkdir dist
	@for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && tar -cvzf ../../dist/$$f.tar.gz *); \
		(cd $(shell pwd)/dist && shasum -a 512 $$f.tar.gz > $$f.sha512); \
		echo $$f; \
	done

release: dist
	@latest_tag=$$(git describe --tags `git rev-list --tags --max-count=1`); \
	comparison="$$latest_tag..HEAD"; \
	if [ -z "$$latest_tag" ]; then comparison=""; fi; \
	changelog=$$(git log $$comparison --oneline --no-merges --reverse); \
	github-release c4milo/$(NAME) $(VERSION) "$$(git rev-parse --abbrev-ref HEAD)" "**Changelog**<br/>$$changelog" 'dist/*'; \
	git pull

deps:
	go get github.com/c4milo/github-release
	go get github.com/mitchellh/gox

updatedeps:
	go list ./... \
        | xargs go list -f '{{join .Deps "\n"}}' \
        | grep -v unicreds \
        | grep -v '/internal/' \
        | sort -u \
        | xargs go get -f -u -v

watch:
	scantest

packages:
	rm -rf package && mkdir package
	rm -rf stage && mkdir -p stage/usr/bin
	cp build/Linux/unicreds stage/usr/bin
	fpm --name $(NAME) -a x86_64 -t rpm -s dir --version $(VERSION) --iteration $(ITERATION) -C stage -p package/$(NAME)-$(VERSION)_$(ITERATION).rpm usr
	fpm --name $(NAME) -a x86_64 -t deb -s dir --version $(VERSION) --iteration $(ITERATION) -C stage -p package/$(NAME)-$(VERSION)_$(ITERATION).deb usr

generate-mocks:
	mockery -dir ../../aws/aws-sdk-go/service/kms/kmsiface --all
	mockery -dir ../../aws/aws-sdk-go/service/dynamodb/dynamodbiface -testonly -all

.PHONY: build fmt test integration watch release packages
