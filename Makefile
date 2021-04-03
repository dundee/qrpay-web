NAME := qrpay-web
PACKAGE := github.com/dundee/$(NAME)
VERSION := $(shell git describe --tags 2>/dev/null)
GOFLAGS ?= -buildmode=pie -trimpath -mod=readonly -modcacherw
LDFLAGS := -s -w -extldflags '-static' \
	-X '$(PACKAGE)/build.Version=$(VERSION)' \
	-X '$(PACKAGE)/build.User=$(shell id -u -n)' \
	-X '$(PACKAGE)/build.Time=$(shell LC_ALL=en_US.UTF-8 date)'

all: clean build-all clean-uncompressed-dist shasums

build:
	@echo "Version: " $(VERSION)
	mkdir -p dist
	GOFLAGS="$(GOFLAGS)" CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o dist/$(NAME) .

build-all:
	@echo "Version: " $(VERSION)
	-mkdir dist
	-CGO_ENABLED=0 gox \
		-os="darwin linux" \
		-arch="amd64" \
		-output="dist/qrpay_{{.OS}}_{{.Arch}}" \
		-ldflags="$(LDFLAGS)"

	cd dist; for file in qrpay_linux_* qrpay_darwin_* ; do tar czf $$file.tgz $$file; done

run:
	go run -tags live .

clean:
	-rm -r dist

clean-uncompressed-dist:
	find dist -type f -not -name '*.tgz' -not -name '*.zip' -delete

shasums:
	cd dist; sha256sum * > sha256sums.txt
	cd dist; gpg --sign --armor --detach-sign sha256sums.txt

.PHONY: build build-all clean
