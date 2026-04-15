.PHONY: gen test utest itest

gen:
	go run github.com/vektra/mockery/v2@latest

test:
	go test -count=1 ./...

utest:
	go test ./... --short

itest:
	go test ./internal/storage/postgres/... -v
