EXEC = reggen

TAR = $(EXEC)-$(VER).tar.gz
BIN = $(DESTDIR)/usr/bin

BUILD_TIME=`date +%H:%M:%S.%Y-%m-%d`
GIT_COMMIT=`git log --pretty=format:"%h" -1`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`

# Add build-time-string into the executable file.
X_ARGS += -X main.BUILD_INFO="$(BUILD_TIME).git:$(GIT_COMMIT)@$(GIT_BRANCH)"

.PHONY: build simple complex install clean archive
.DEFAULT_GOAL := help

build: ## build the target binary.
	@go build -ldflags "$(X_ARGS)" -o $(EXEC)

test:  ## run unit test.
	@go test ./...

simple: build ## build target and run with the simple.regs example.
	@./$(EXEC) -i testdata/simple.regs -d | less

complex: build ## build target and run with the complex.regs example.
	@./$(EXEC) -i testdata/complex.regs -d | less

install:$(EXEC) ## install the target binary.
	install -d $(BIN)
	install $(EXEC) $(BIN)

clean: ## clean the target binary.
	@rm -rfv $(EXEC)

archive: ## archive source code to a tar ball with version information.
	@echo "archive to $(TAR)"
	@git archive master --prefix="$(EXEC)-$(VER)/" --format tar.gz -o $(TAR)

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'
