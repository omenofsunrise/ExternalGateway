.PHONY: proto
proto:
	@echo "Generating proto files..."
	buf generate
	@echo "✅ Proto files generated"

.PHONY: build
build:
	go build -o bin/server cmd/server/main.go

.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: clean
clean:
	rm -rf api/gateway/
	rm -rf bin/

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: all
all: proto tidy build