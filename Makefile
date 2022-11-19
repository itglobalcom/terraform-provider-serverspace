HOSTNAME=serverspace.by
NAMESPACE=main
NAME=serverspace
BINARY=terraform-provider-${NAME}
VERSION=0.2
OS_ARCH=linux_amd64

default: install

build:
	go build -o ${BINARY}

release: build_darwin_amd64 build_freebsd_386 build_freebsd_amd64 \
		 build_freebsd_arm build_linux_386 build_linux_amd64 \
		 build_linux_arm build_openbsd_386 build_openbsd_amd64 \
		 build_solaris_amd64 build_windows_386 build_windows_amd64


build_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64

build_freebsd_386:
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386

build_freebsd_amd64:
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64

build_freebsd_arm:
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm

build_linux_386:
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386

build_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64

build_linux_arm:
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm

build_openbsd_386:
	GOOS=openbsd GOARCH=386 go build -o ./bin/${NARY}_${VERSION}_openbsd_386

build_openbsd_amd64:
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64

build_solaris_amd64:
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64

build_windows_386:
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386

build_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64


install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
