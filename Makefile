build:
	@go build -o bin/bankWithGo

run: build 
	@./bin/bankWithGo

test:
	@go test -v ./..