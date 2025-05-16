vendor:
	@echo "Running go mod vendor..."
	@go mod vendor

run:
	@echo "Running app..."
	@go run cmd/gql/main.go

generate:
	@echo "Running generate file..."
	@go generate ./...

gql-gen:
	@echo "Running gql generator..."
	@go run github.com/99designs/gqlgen generate

test:
	@echo "Running tests..."
	@go test -v $$(go list ./... | grep -v /vendor/ | grep -v /graph/ | grep -v /cmd/)