
.PHONY: gen
gen:
	go run github.com/vektra/mockery/v2@latest

.PHONY: test
test:
	go test -count=1 ./...

.PHONY: utest
utest:
	go test -count=1 ./... --short