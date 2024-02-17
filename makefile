all:
	mkdir -p build
	CGO_ENABLED=1 CC=gcc GOARCH=amd64 go build -tags static -ldflags "-s -w" -o build/gb.exe

clean:
	rm -rf build

run:
	go run main.go