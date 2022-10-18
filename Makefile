build-test:
	@go build -o out/busy_lock ./tests/busy_lock
	@go build -o out/crash_lock ./tests/crash_lock
