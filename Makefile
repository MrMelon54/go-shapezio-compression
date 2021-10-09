

all:
	make build-main

build-main:
	mkdir -p dist/ && \
	go build -o dist/shapezio-compression ./cmd/shapezio-compression

package:
	./scripts/package-shapezio-compression.sh

deb:
	make clean && \
	make all && \
	make package

clean:
	rm -rf dist/ && \
	rm -rf package/


.PHONY: run
run:
	cd ./dist/ && \
	./shapezio-compression
