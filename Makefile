.PHONY: all
all: build run

build:
	docker-compose build
	# The command below builds only a go application image.
	# docker build -t t2-http .

run:
	docker-compose up
	# The command below runs only a go application container.
	# docker run -it --rm --name t2-http t2-http

test:
	go test ./...

gen:
	go generate ./...
