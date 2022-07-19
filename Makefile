CWD_ABS				= $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

NAME				= my-service

.phony: lint init build test
.DEFAULT_GOAL = $(NAME)

init:
	if [ ! -f "go.mod" ]; then\
		go mod init $(NAME);\
		go mod tidy;\
	fi
	go get $(CWD_ABS)...

build:
	go build $(CWD_ABS)...

test:
	go test -v --cover $(CWD_ABS)...

lint:
	echo 'Mock for linters...'

clean:
	go clean --testcache
	go mod tidy

clear: clean
	if [ -f "${CWD_ABS}${NAME}" ]; then\
		rm -v "${CWD_ABS}${NAME}";\
	fi

$(NAME): build

