APP_NAME = bugbox
GO = go
GOFMT = gofmt

all: build

build:
	@echo "Building..."
	$(GO) build -o $(APP_NAME) main.go

run: build
	@echo "Running the application..."
	./$(APP_NAME)

watch:
	@echo "Watching for changes..."
	nodemon --watch './**/*.go' --ext go --signal SIGTERM --exec 'make run bugbox || exit 1'

fmt:
	@echo "Formatting Go code..."
	$(GOFMT) -w .

test:
	@echo "Running tests..."
	$(GO) test -v ./...

clean:
	@echo "Cleaning build..."
	rm -f $(APP_NAME)

install:
	@echo "Installing dependencies..."
	$(GO) mod tidy

release:
	@echo "Releasing..."
	goreleaser release --clean
