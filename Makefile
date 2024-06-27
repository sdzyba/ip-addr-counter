.PHONY: run
# Execute ip-addr-counter with a given f (filepath) and c (concurrency) parameters
# @Example
# 	$ make run f=test c=500
run:
	$(eval FILEPATH = $(if ${f},${f},ip_addresses))
	$(eval CONCURRENCY = $(if ${c},${c},10000))
	@go run . ${FILEPATH} ${CONCURRENCY}

.PHONY: test
test:
	@go test -race -v .

.PHONY: benchmark
benchmark:
	@go test -bench=.

.PHONY: start-pprof-allocs-web
start-pprof-allocs-web:
	@go tool pprof -http=:8080 http://localhost:6060/debug/pprof/allocs

.PHONY: start-pprof-heap-web
start-pprof-heap-web:
	@go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap
