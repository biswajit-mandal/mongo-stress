build:
	go build
deps:
	go get -v ./...
stress: build
	bash stress-test.sh
