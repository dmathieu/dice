GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

build: gox
	gox -osarch="linux/amd64 linux/arm linux/arm64" -output="compiled/{{.Dir}}_{{.OS}}_{{.Arch}}"

ci: tidy test

test:
	go test -mod=vendor -race -v -coverprofile c.out ./...

tidy: goimports
	test -z "$$(goimports -l -d $(GO_FILES) | tee /dev/stderr)"
	go vet -mod=vendor ./...

goimports:
	go get golang.org/x/tools/cmd/goimports

gox:
	go get github.com/mitchellh/gox
