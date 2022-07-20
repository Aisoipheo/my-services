CWD_ABS				= $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

NAME				= my-service

.phony: lint init build test coverage
.DEFAULT_GOAL = $(NAME)

init:
	if [ ! -f "go.mod" ]; then\
		go mod init $(NAME);\
		go mod tidy;\
	fi
	go get ./...

build:
	go build ./...

test:
	go test -v --cover ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	./bin/golangci-lint run ./...

clean:
	rm -vf coverage.out
	go clean --testcache
	go mod tidy

clear: clean
	rm -vf "${NAME}"

$(NAME): build

