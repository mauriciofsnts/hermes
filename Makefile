
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

# auto restart
.PHONY: dev
dev:
	air

.PHONY: lint
lint:
	revive -formatter friendly -config revive.toml ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: gosec
gosec:
	gosec -tests ./... 

.PHONY: inspect
inspect: lint gosec staticcheck

.PHONY: install-inspect-tools
install-inspect-tools:
	go install github.com/mgechev/revive@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest