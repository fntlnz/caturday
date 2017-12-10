.PHONY: build clean

build:
	docker build -t docker.io/fntlnz/caturday .

push:
	docker push docker.io/fntlnz/caturday
