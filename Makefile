all: update-deps build
build:
	go build -mod vendor -o github-clone .
.PHONY: build

debug:
	go build -gcflags="all=-N -l"  -mod vendor -o github-clone .
.PHONY: build

update-deps:
	GO111MODULE=on go mod vendor
.PHONY: update-deps

run: check-env
	./github-clone
.PHONY: run

clean:
	rm -rf ./github-clone
.PHONY: clean

check-env:
ifndef GITHUB_OAUTH_TOKEN
	$(error GITHUB_OAUTH_TOKEN is not set)
endif


