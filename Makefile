PROJECT		:=go-demo
MAIN_DIR	:=./cmd/$(PROJECT)
SOURCES		:=$(shell find . -type f -name '*.go')
PORT 		:=$(shell grep 'HttpPort' conf/app.ini | cut -d '=' -f2 | tr -d ' ')
BASE_URL	:=http://localhost:$(PORT)

define HELP_TEXT
release: Build docker container
build: Build local
run: Run $(PROJECT) at local
test: Run test cases
start: Start docker compose
redis: Start docker compose only Redis, before you run applicaion at local
clean: Remove object files and cached files
format: Format sources
docs: Gernerate /docs
create: curl test to create task by picking random word
update: curl test to update task by example 'make update id=10 text=movie status=1'
delete: curl test to delete task by example 'make delete id=10'
list: curl test to list tasks
endef
export HELP_TEXT

.PHONY: all help release build run start redis clean format gotestsum FORCE create update list delete swag doc

all: help

help:
	@echo "$(PROJECT) - Makefile commands"
	@echo
	@echo "make [options]"
	@echo
	@echo "$$HELP_TEXT" | while IFS=: read -r command description; do \
		printf "%-10s %s\n" "$$command" "$$description"; \
	done

release: .docker_build

.docker_build: $(SOURCES) ./Dockerfile docs
	@docker build --progress=plain -t ${PROJECT} .
	@touch $@

$(PROJECT): $(SOURCES) FORCE
	@go build -v $(MAIN_DIR)
	@echo "Build finished successfully."

build: $(PROJECT)

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

docs: swag
	@swag init -g $(MAIN_DIR)/main.go 

swag:
	@command -v swag >/dev/null 2>&1 || {\
		echo >&2 "swag is not installed. Installing..."; \
		GO111MODULE=on go install github.com/swaggo/swag/cmd/swag@v1.8.1; \
	}

list:
	@TMPFILE=$$(mktemp); \
	STATUS_CODE=$$(curl -s -o $$TMPFILE -w "%{http_code}" -X GET $(BASE_URL)/tasks); \
	echo "$$STATUS_CODE GET $(BASE_URL)/tasks"; \
	if [ $$STATUS_CODE -ge 200 ] && [ $$STATUS_CODE -lt 400 ]; then \
		cat $$TMPFILE | python3 -m json.tool; \
	else \
		echo "Request failed with status code: $$STATUS_CODE"; \
	fi; \
	rm -f $$TMPFILE

create:
	@TMPFILE=$$(mktemp); \
	STATUS_CODE=$$(curl -s -o $$TMPFILE -w "%{http_code}" -X POST -H "Content-Type: application/json" -d '{"text": "'$(shell sort -R /usr/share/dict/words | head -n 1)'"}' $(BASE_URL)/task); \
	echo "$$STATUS_CODE POST $(BASE_URL)/task"; \
	if [ $$STATUS_CODE -ge 200 ] && [ $$STATUS_CODE -lt 400 ]; then \
		cat $$TMPFILE | python3 -m json.tool; \
	else \
		echo "Request failed with status code: $$STATUS_CODE"; \
	fi; \
	rm -f $$TMPFILE

update:
	@if [ -z "$(id)" ]; then \
		echo "Update which task?"; \
		echo 'Use "make edit id=10 text=ski status=1"'; \
		exit 1; \
	fi; \
	if [ -z "$(text)" -a -z "$(status)" ]; then \
		echo "Update what field?"; \
		echo 'Use "make edit id=10 text=ski status=1"'; \
		exit 1; \
	fi;
	$(eval TMPFILE := $(shell mktemp))
	@BODY="{\"id\": $(id)"; \
	if [ -n "$(text)" ]; then BODY="$$BODY, \"text\": \"$(text)\""; fi; \
	if [ -n "$(status)" ]; then BODY="$$BODY, \"status\": $(status)"; fi; \
	BODY="$$BODY }"; \
	STATUS_CODE=$$(curl -s -o $(TMPFILE) -w "%{http_code}" -X PUT -H "Content-Type: application/json" -d "$$BODY" $(BASE_URL)/task/$(id)); \
	echo "$$STATUS_CODE PUT $(BASE_URL)/task/$(id)"; \
	if [ $$STATUS_CODE -ge 200 ] && [ $$STATUS_CODE -lt 400 ]; then \
		cat $(TMPFILE) |  python3 -m json.tool; \
	else \
		cat $(TMPFILE); \
	fi; \
	echo; \
	rm -f $(TMPFILE)

delete:
	@if [ -z "$(id)" ]; then \
		echo "Delete which task?"; \
		echo 'Use "make delete id=10"'; \
		exit 1; \
	fi;
	@TMPFILE=$$(mktemp); \
	STATUS_CODE=$$(curl -s -o $$TMPFILE -w "%{http_code}" -X DELETE $(BASE_URL)/task/$$id); \
	echo "$$STATUS_CODE DELETE $(BASE_URL)/task/$(id)"; \
	if [ $$STATUS_CODE -ge 400 ] ; then \
		cat $$TMPFILE; \
	fi; \
	rm -f $$TMPFILE
