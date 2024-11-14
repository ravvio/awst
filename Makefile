
all: build

build:
	@echo "Building..."
	@go build -o awst main.go

run:
	@go run main.go

test:
	@echo "Testing..."
	@go test ./tests -v

clean:
	@echo "Cleaning..."
	@rm -f main
