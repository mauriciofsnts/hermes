
BINARY_NAME = hermes
TEST_COMMAND = go test

.PHONY: build
build:
	go build -v -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

.PHONY: dist
dist:
	CGO_ENABLED=0 go build -gcflags=all=-l -v -ldflags="-w -s" -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

.PHONY: run
run: build
	./$(BINARY_NAME)

.PHONY: start
start: build setup-config
	./$(BINARY_NAME)

.PHONY: setup-config
setup-config:
	@if [ ! -f "config.yaml" ]; then \
		echo "📝 Creating config.yaml from config_example.yaml..."; \
		cp config_example.yaml config.yaml; \
		echo "✅ config.yaml created"; \
	fi

.PHONY: test
test:
	$(TEST_COMMAND) -cover -parallel 5 -failfast  ./...

.PHONY: test-integration
test-integration:
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker compose -f docker-compose.test.yml down

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: swagger
swagger:
	swag init -g cmd/hermes/main.go -o docs --parseDependency --parseInternal

# auto restart
.PHONY: dev
dev:
	go tool air

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
