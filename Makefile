.PHONY: clean
clean:
	make -C service clean

.PHONY: gen
gen:
	make -C api gen
	make -C service gen

.PHONY: test
test:
	make -C service test

.PHONY: build
build:
	make -C service build

.PHONY: lint
lint:
	make -C service lint

.PHONY: run
run:
	make -C service run

.PHONY: test-integration
test-integration:
	make -C service test-integration

.PHONY: test-integration-race
test-integration-race:
	make -C service integration-test-race

install-tools:
	go install github.com/golang/mock/mockgen@v1.6.0;
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0;
	brew install openapi-generator;
