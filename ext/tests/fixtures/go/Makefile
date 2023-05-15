all: test.so

test.so: test.go ../../../../php-go/*
	go build -o test.so -buildmode c-shared
