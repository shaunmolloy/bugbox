APP_NAME = bugbox
GO = go
GOFMT = gofmt
EXCLUDE_DIRS = /cmd /internal/scheduler /internal/tui
INCLUDE_DIRS=$(shell \
	go list ./... | \
	grep -v -F $(foreach p,$(EXCLUDE_DIRS),-e $(p)) \
)

all: build

build:
	@echo "Building..."
	$(GO) build -o $(APP_NAME) main.go

run: build
	@echo "Running the application..."
	./$(APP_NAME)

fmt:
	@echo "Formatting Go code..."
	$(GOFMT) -w .

test:
	@echo "Running tests..."
	$(GO) test $(INCLUDE_DIRS) -v ./...

test_cover:
	@echo "Running tests with coverage..."
	$(GO) test $(INCLUDE_DIRS) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	$(GO) tool cover -func coverage.out
	@echo "Coverage report generated: coverage.html"

test_watch:
	@echo "Watching for changes..."
	nodemon --watch './**/*.go' --ext go --signal SIGTERM --exec 'make test_cover || exit 1'

clean:
	@echo "Cleaning build..."
	rm -f $(APP_NAME)

install:
	@echo "Installing dependencies..."
	$(GO) mod tidy
