.PHONY: compile deps fmt test


all: compile


compile: deps
	@go build

fmt:
	@go fmt

test:
	@go test ./...

clean:
	@go clean
	@rm -rf _vendor
