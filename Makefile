TARGET=remove-slowly

ALL: build

build: dist */*.go
	env CGO_ENABLED=0 go build -o dist/remove-slowly ./cli

release-local-dryrun:
	goreleaser --snapshot --skip-publish --clean

dist: $@
	mkdir $@

.PHONY: clean test build

test: clean
	@mkdir -p test /tmp/test_results cli/test/foo cli/test/bar
	touch cli/test/file0 cli/test/foo/file1 cli/test/foo/file2
	gotestsum --junitfile /tmp/test_results/unit-tests.xml -- -coverprofile=./test/coverage.out ./...
	go tool cover -html=test/coverage.out -o test/coverage.html

clean:
	- $(RM) -rf dist/* test/* cli/test
