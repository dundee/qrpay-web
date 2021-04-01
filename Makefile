NAME := qrcode-payment-web
PACKAGE := github.com/dundee/$(NAME)
VERSION := $(shell git describe --tags 2>/dev/null)
GOFLAGS ?= -buildmode=pie -trimpath -mod=readonly -modcacherw
LDFLAGS := -s -w -extldflags '-static' \
	-X '$(PACKAGE)/build.Version=$(VERSION)' \
	-X '$(PACKAGE)/build.User=$(shell id -u -n)' \
	-X '$(PACKAGE)/build.Time=$(shell LC_ALL=en_US.UTF-8 date)'

build:
	@echo "Version: " $(VERSION)
	mkdir -p dist
	GOFLAGS="$(GOFLAGS)" CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o dist/$(NAME) .
