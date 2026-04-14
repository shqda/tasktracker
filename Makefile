.PHONY: gen test utest

gen:
	go run github.com/vektra/mockery/v2@latest

utest:
	go test ./... --short

test:
	go test -count=1 ./...

