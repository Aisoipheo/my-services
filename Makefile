CWD_ABS				= $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

NAME				= feed-service

.phony: lint init build test coverage bin_dir
.DEFAULT_GOAL = $(NAME)

bin_dir:
	if [ ! -d "./bin" ]; then\
		mkdir "bin";\
	fi

init:
	if [ ! -f "go.mod" ]; then\
		go mod init $(NAME);\
		go mod tidy;\
	fi
	go get ./...

build: init bin_dir
	go build -o ./bin/$(NAME) ./cmd/$(NAME)

test:
	go test -cpu=4 -v --cover ./...

coverage:
	# ignore errors, generate html for valid coverage
	-go test -cpu=4 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./...
	yamllint ./...
	hadolint ./...

clean:
	rm -vf coverage.out
	go clean --testcache
	go mod tidy

fclean: clean
	rm -vf ./bin/$(NAME) go.mod go.sum

$(NAME): build

