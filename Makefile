
GO=go1.18beta1

test:
	${GO} test -fuzztime=1m -fuzz ./...

.PHONY: test
