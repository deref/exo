PACKAGE_NAME          := github.com/deref/exo
GOLANG_CROSS_VERSION  ?= v1.17.2

.PHONY: all
all:
	$(MAKE) codegen
	$(MAKE) build

.PHONY: build
build: bin/exo
	$(MAKE) -C gui build

.PHONY:
bin/exo:
	go build -o ./bin/exo

.PHONY: codegen
codegen:
	./script/codegen.sh

.PHONY: make-gui
make-gui:
	make -C gui

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: run-tests
run-tests: bin/exo
	go run ./test/main.go ./bin/exo ./test/image/fixtures

.PHONY: release-dry-run
release-dry-run: make-gui mod-tidy codegen
	@docker run \
		--privileged \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/troian/golang-cross:${GOLANG_CROSS_VERSION} \
		--rm-dist --skip-validate --skip-publish

.PHONY: release
release: make-gui mod-tidy codegen
	@docker run \
		--privileged \
		-e GITHUB_TOKEN \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/troian/golang-cross:${GOLANG_CROSS_VERSION} \
		--rm-dist
