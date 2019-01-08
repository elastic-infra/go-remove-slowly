TARGET=remove-slowly

ALL: build

build: dist */*.go
	env CGO_ENABLED=0 gox -os "linux" -arch "amd64 386 arm arm64" -output dist/remove-slowly_{{.OS}}_{{.Arch}} ./...

dist: $@
	mkdir $@

.PHONY: clean test build

test:
	@mkdir -p test /tmp/test_results
	gotestsum --junitfile /tmp/test_results/unit-tests.xml -- -coverprofile=./test/coverage.out ./...
	go tool cover -html=test/coverage.out -o test/coverage.html

clean:
	- $(RM) dist/* test/*
