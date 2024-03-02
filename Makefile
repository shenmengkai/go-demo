PROJECT		:=gogolook2024
MAIN_DIR	:=./cmd/$(PROJECT)

define HELP_TEXT
release: Build docker container
build: Build local
run: Run $(PROJECT) at local
test: Run test cases
start: Start docker compose
redis: Start docker compose only Redis
clean: Remove object files and cached files
format: Format sources
endef
export HELP_TEXT

.PHONY: release build run help start redis clean format gotestsum

all: help

help:
	@echo "$(PROJECT) - Makefile commands"
	@echo
	@echo "make [options]"
	@echo
	@echo "$$HELP_TEXT" | while IFS=: read -r command description; do \
		printf "%-10s %s\n" "$$command" "$$description"; \
	done

release: build
	@docker build --progress=plain -t ${PROJECT} .

build:
	@go build -v $(MAIN_DIR)
	@echo Build $(MAIN_DIR)

clean:
	@rm -rf ${PROJECT}
	@go clean -x -r -cache -testcache -modcache -fuzzcache .
	@echo Clean $(MAIN_DIR)

run:
	@go run $(MAIN_DIR)

test: build gotestsum
	@gotestsum --format testname ./internal/middleware ./internal/service

start: release
	@docker compose up

redis:
	@docker compose up redis

format:
	go vet ./...; true
	gofmt -w .

gotestsum:
	@command -v gotestsum >/dev/null 2>&1 || {\
		echo >&2 "gotestsum is not installed. Installing..."; \
		GO111MODULE=on go install gotest.tools/gotestsum@latest; \
	}

list:
	@clear
	@curl -X GET http://localhost:8000/tasks

add:
	@clear
	@curl -X POST http://localhost:8000/task -H "Content-Type: application/json" -d '{"text": "'$(shell sort -R /usr/share/dict/words | head -n 1)'"}'

edit:
	@clear
	@curl -X PUT http://localhost:8000/task/17 -H "Content-Type: application/json" -d '{"text": "rabbit",  "status":0}'

del:
	@clear
	@curl -X DELETE http://localhost:8000/task/17
