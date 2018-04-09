EXEC = reggen

TAR = $(EXEC)-$(VER).tar.gz
BIN = $(DESTDIR)/usr/bin

all:bin

bin:
	@go build -o $(EXEC)

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
