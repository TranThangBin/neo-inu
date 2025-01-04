all: build

SCRIPT := ./script/main.sh

build_bin:
	go build -o ./bin/neo-inu ./cmd/neo-inu/main.go

make clean_bin:
	rm -rf ./bin ./tmp

tag:
	$(SCRIPT) tag

build:
	$(SCRIPT) build

run:
	$(SCRIPT) run

restart:
	$(SCRIPT) restart

stop:
	$(SCRIPT) stop

clean_container:
	$(SCRIPT) clean_container

clean_image:
	$(SCRIPT) clean_image

clean: stop clean_container clean_image

.PHONY: build
