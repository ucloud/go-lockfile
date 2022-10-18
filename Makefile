.PHONY: build-test
build-test:
	@go build -o out/busy_lock ./example/busy_lock
	@go build -o out/crash_lock ./example/crash_lock
	@go build -o out/sample_lock ./example/sample_lock
	@go build -o out/competition ./example/competition

.PHONY: test-competition
test-competition: build-test
	@./out/competition 10

.PHONY: fmt
fmt:
	@find . -name \*.go -exec goimports -w {} \;
