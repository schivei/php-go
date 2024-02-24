all: test.so

test.so: test.go ../../../../php/*
	CGO_ENABLED=1 go build -o test.so -buildmode c-shared
