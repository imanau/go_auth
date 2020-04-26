NAME=go_auth
VERSION=1.0

MAKEFILE_DIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))

up:
	docker-compose up -d db db_test &&   \
	sleep 3 && \
	docker-compose up go

build:
	docker build -t $(NAME) .

bash:
	docker run -it -v $(MAKEFILE_DIR):/usr/src/go_auth -p 3000:3000  $(NAME)

run:
	docker run -it -v $(MAKEFILE_DIR):/usr/src/go_auth -p 3000:3000  $(NAME) bash -c 'go run main.go'

test:
	docker-compose exec go /bin/bash -c 'GO_ENV=test go test -v ./...'
