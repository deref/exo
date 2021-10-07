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
