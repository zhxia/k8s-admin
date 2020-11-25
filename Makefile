### --- Makefile ---

progname=kube-admin
all: build
build:export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig
build:export CGO_LDFLAGS=-L/usr/local/lib -Wl,-rpath -Wl,$$ORIGIN/lib
build:export GOPROXY=https://goproxy.io
build:
	mkdir -p ./bin
	cp config.yml ./bin
	go build -ldflags "-X 'main.version=2.0.0'" -o bin/$(progname) src/main.go
linux:export CGO_ENABLED=0
linux:export GOARCH=amd64
linux:export GOOS=linux
linux:
	go build -ldflags "-X 'main.version=2.0.0'" -o bin/$(progname) src/main.go
clean:
	rm -rf ./bin