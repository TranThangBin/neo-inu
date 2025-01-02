all: build

SCRIPT := ./script/main.sh

build_bin:
	go build -o ./bin/neo-inu ./cmd/neo-inu/main.go

build:
	$(SCRIPT) build

run:
	$(SCRIPT) run

stop:
	$(SCRIPT) stop

restart:
	$(SCRIPT) restart

clean_container:
	$(SCRIPT) clean_container

clean_image:
	$(SCRIPT) clean_image

clean: stop clean_container clean_image

.PHONY: build
