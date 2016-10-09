.PHONY: build clean

build:
	mkdir -p dist
	docker build -t fntlnz/caturday-build -f Dockerfile.build .
	docker run --rm fntlnz/caturday-build cat caturday > dist/caturday
	chmod +x dist/caturday
	docker build -t fntlnz/caturday .

clean:
	rm -Rf dist/
