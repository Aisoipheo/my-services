CWD_ABS				= $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

NAME				= my-service

.phony: lint init build
.DEFAULT_GOAL = $(NAME)

init:
	go mod init $(NAME)
	go mod tidy

build:
	go build $(CWD_ABS)...

test:
	go test -cover $(CWD_ABS)... -v -coverpkg=$(CWD_ABS)...

lint:
	echo 'Mock for linters...'

$(NAME): build

