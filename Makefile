
name = awst

all: build

build:
	@echo "Building..."
	@go build -o ${name} main.go

run:
	@go run main.go

test:
	@echo "Testing..."
	@go test ./tests -v

clean:
	@echo "Cleaning..."
	@rm -f ${name}
	@go mod tidy
	@go clean
