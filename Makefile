CWD_ABS				= $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

NAME				= feed-service

.phony: lint init build test coverage bin_dir golangci_lint_install
.DEFAULT_GOAL = $(NAME)

bin_dir:
	if [ ! -d "./bin" ]; then\
		mkdir "bin";\
	fi

golangci_lint_install:
	if [ ! -f "./bin/golangci-lint" ]; then\
		wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.47.1;\
	fi

init:
	if [ ! -f "go.mod" ]; then\
		go mod init $(NAME);\
		go mod tidy;\
	fi
	go get ./...

build: bin_dir
	go build -o ./bin/$(NAME) ./cmd/$(NAME)

test:
	go test -v --cover ./...

coverage:
	# ignore errors, generate html for valid coverage
	-go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint: bin_dir golangci_lint_install
	./bin/golangci-lint run ./...

clean:
	rm -vf coverage.out
	go clean --testcache
	go mod tidy

clear: clean
	rm -vf "${NAME}"

$(NAME): build

