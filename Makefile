EXEC = reggen

TAR = $(EXEC)-$(VER).tar.gz
BIN = $(DESTDIR)/usr/bin

BUILD_TIME=`date +%H:%M:%S.%Y-%m-%d`
GIT_COMMIT=`git log --pretty=format:"%h" -1`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`

# Add build-time-string into the executable file.
X_ARGS += -X main.BUILD_INFO="$(BUILD_TIME).git:$(GIT_COMMIT)@$(GIT_BRANCH)"

all:bin

bin:
	@echo "$(BUILD_TIME) $(GIT_COMMIT)@$(GIT_BRANCH)"
	@go build -ldflags "$(X_ARGS)" -o $(EXEC)

install:$(EXEC)
	install -d $(BIN)
	install $(EXEC) $(BIN)

clean:
	@rm -rfv $(EXEC)

archive:
	@echo "archive to $(TAR)"
	@git archive master --prefix="$(EXEC)-$(VER)/" --format tar.gz -o $(TAR)

test:
	@go test
