PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
LDFLAGS := $(shell go run buildscripts/gen-ldflags.go)

BUILD_LDFLAGS := '$(LDFLAGS)'

all: build

checks:
	@echo "Checking dependencies"
	@(env bash $(PWD)/buildscripts/checkdeps.sh)
	@echo "Checking for project in GOPATH"
	@(env bash $(PWD)/buildscripts/checkgopath.sh)

getdeps:
	@echo "Installing golint" && go get -u golang.org/x/lint/golint
	@echo "Installing gocyclo" && go get -u github.com/fzipp/gocyclo
	@echo "Installing deadcode" && go get -u github.com/remyoudompheng/go-misc/deadcode
	@echo "Installing misspell" && go get -u github.com/client9/misspell/cmd/misspell
	@echo "Installing ineffassign" && go get -u github.com/gordonklaus/ineffassign

verifiers: getdeps vet fmt lint cyclo deadcode spelling

vet:
	@echo "Running $@"
	@go tool vet -atomic -bool -copylocks -nilfunc -printf -shadow -rangeloops -unreachable -unsafeptr -unusedresult cmd
	@go tool vet -atomic -bool -copylocks -nilfunc -printf -shadow -rangeloops -unreachable -unsafeptr -unusedresult pkg

fmt:
	@echo "Running $@"
	@gofmt -d cmd
	@gofmt -d pkg

lint:
	@echo "Running $@"
	@${GOPATH}/bin/golint -set_exit_status github.com/didiyun/mc/cmd...
	@${GOPATH}/bin/golint -set_exit_status github.com/didiyun/mc/pkg...

ineffassign:
	@echo "Running $@"
	@${GOPATH}/bin/ineffassign .

cyclo:
	@echo "Running $@"
	@${GOPATH}/bin/gocyclo -over 100 cmd
	@${GOPATH}/bin/gocyclo -over 100 pkg

deadcode:
	@echo "Running $@"
	@${GOPATH}/bin/deadcode -test $(shell go list ./...)

spelling:
	@${GOPATH}/bin/misspell -error `find cmd/`
	@${GOPATH}/bin/misspell -error `find pkg/`
	@${GOPATH}/bin/misspell -error `find docs/`

# Builds minio, runs the verifiers then runs the tests.
check: test
test: verifiers build
	@echo "Running unit tests"
	@go test $(GOFLAGS) -tags kqueue ./...
	@echo "Running functional tests"
	@(env bash $(PWD)/functional-tests.sh)

coverage: build
	@echo "Running all coverage for minio"
	@(env bash $(PWD)/buildscripts/go-coverage.sh)

# Builds minio locally.
build: checks
	@echo "Building minio binary to './mc'"
	@CGO_ENABLED=0 go build -tags kqueue --ldflags $(BUILD_LDFLAGS) -o $(PWD)/mc

pkg-add:
	@echo "Adding new package $(PKG)"
	@${GOPATH}/bin/govendor add $(PKG)

pkg-update:
	@echo "Updating new package $(PKG)"
	@${GOPATH}/bin/govendor update $(PKG)

pkg-remove:
	@echo "Remove new package $(PKG)"
	@${GOPATH}/bin/govendor remove $(PKG)

pkg-list:
	@$(GOPATH)/bin/govendor list

# Builds minio and installs it to $GOPATH/bin.
install: build
	@echo "Installing mc binary to '$(GOPATH)/bin/mc'"
	@mkdir -p $(GOPATH)/bin && cp $(PWD)/mc $(GOPATH)/bin/mc
	@echo "Installation successful. To learn more, try \"mc --help\"."

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@rm -rvf mc
	@rm -rvf build
	@rm -rvf release


build = GOOS=$(1) GOARCH=$(2) go build -tags kqueue --ldflags $(BUILD_LDFLAGS) -o build/mc$(3)
tar = cd build && tar -cvzf $(1)_$(2).tar.gz mc$(3) && rm mc$(3)
zip = cd build && zip $(1)_$(2).zip mc$(3) && rm mc$(3)

release: windows_build darwin_build linux_build bsd_build

##### WINDOWS BUILDS #####
windows_build: build/windows_386.zip build/windows_amd64.zip

build/windows_386.zip: checks
	$(call build,windows,386,.exe)
	$(call zip,windows,386,.exe)

build/windows_amd64.zip: checks
	$(call build,windows,amd64,.exe)
	$(call zip,windows,amd64,.exe)

##### LINUX BUILDS #####
linux_build: build/linux_arm.tar.gz build/linux_arm64.tar.gz build/linux_386.tar.gz build/linux_amd64.tar.gz

build/linux_386.tar.gz: checks
	$(call build,linux,386,)
	$(call tar,linux,386)

build/linux_amd64.tar.gz: checks
	$(call build,linux,amd64,)
	$(call tar,linux,amd64)

build/linux_arm.tar.gz: checks
	$(call build,linux,arm,)
	$(call tar,linux,arm)

build/linux_arm64.tar.gz:checks
	$(call build,linux,arm64,)
	$(call tar,linux,arm64)

##### DARWIN (MAC) BUILDS #####
darwin_build: build/darwin_amd64.tar.gz

build/darwin_amd64.tar.gz: checks
	$(call build,darwin,amd64,)
	$(call tar,darwin,amd64)

##### BSD BUILDS #####
bsd_build: build/freebsd_arm.tar.gz build/freebsd_386.tar.gz build/freebsd_amd64.tar.gz \
 build/openbsd_arm.tar.gz build/openbsd_386.tar.gz build/openbsd_amd64.tar.gz

build/freebsd_386.tar.gz:checks
	$(call build,freebsd,386,)
	$(call tar,freebsd,386)

build/freebsd_amd64.tar.gz: checks
	$(call build,freebsd,amd64,)
	$(call tar,freebsd,amd64)

build/freebsd_arm.tar.gz: checks
	$(call build,freebsd,arm,)
	$(call tar,freebsd,arm)

build/openbsd_386.tar.gz: checks
	$(call build,openbsd,386,)
	$(call tar,openbsd,386)

build/openbsd_amd64.tar.gz: checks
	$(call build,openbsd,amd64,)
	$(call tar,openbsd,amd64)

build/openbsd_arm.tar.gz: checks
	$(call build,openbsd,arm,)
	$(call tar,openbsd,arm)
