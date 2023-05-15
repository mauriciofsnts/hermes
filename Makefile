
BINARY_NAME = hermes
TEST_COMMAND = gotest

.PHONY: build
build:
	go build -v -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

.PHONY: dist
dist: 
	CGO_ENABLED=0 go build -gcflags=all=-l -v -ldflags="-w -s" -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

.PHONY: run
run: build
	./$(BINARY_NAME) 

.PHONY: test
test: 
	$(TEST_COMMAND) -cover -parallel 5 -failfast  ./... 

.PHONY: tidy
tidy:
	go mod tidy

# auto restart bot (using fiber CLI <3)
.PHONY: dev
dev:
	fiber dev -t ./cmd/$(BINARY_NAME)