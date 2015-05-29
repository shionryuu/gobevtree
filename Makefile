.PHONY: compile deps fmt test


all: compile


compile: deps
	@gom build ./...

deps: gomfile
	@test -d _vendor || gom install

gomfile:
	@test -f Gomfile || gom gen	

fmt:
	@gom exec go fmt

test:
	@gom test ./...

clean:
	@gom exec go clean
	@rm -rf _vendor
