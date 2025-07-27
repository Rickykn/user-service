.PHONY: setup run all

# Run go mod tidy and vendor dependencies
setup:
	go mod tidy && go mod vendor

# Run your application
run:
	go run main.go

# Combine setup and run
all: setup run
