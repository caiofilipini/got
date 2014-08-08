build:
	go build

test:
	env.sh go test -v

build-linux: clean
	GOOS=linux GOARCH=amd64 go build -o release/got

build-mac: clean
	GOOS=darwin GOARCH=amd64 go build -o release/got

clean:
	rm -fr release
	mkdir release
