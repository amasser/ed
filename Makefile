.PHONY: dev build install release clean

CGO_ENABLED=0

all: dev

dev: build
	@./ed

build: clean
	@go build \
		-tags "netgo static_build" -installsuffix netgo \
		.

install: build
	@go install

release:
	@./tools/release.sh

clean:
	@git clean -f -d -X
