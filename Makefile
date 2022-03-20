get-deps:
	go get -v -t -d .
build:
	go build -o ./dist/app -v ./cmd/
install:
	cp ./dist/app /usr/local/bin/app
