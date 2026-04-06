.PHONY: default target
all:

GO := go
GOTEST := go test -race -v ./...
GOCOVER := golangci-lint run ./...

COVER:
	mkdir -p coverages
	go test -race -coverprofile=coverages/coverage.out ./internal/... ./pkg/...

.PHONY: clean-coverage
	rm -f coverages/coverage.out

.PHONY: report
	go tool-cover -func=cover ./... > coverages/coverage.out 2>&1
	@echo "Coverage Report:"
	@echo "=================="
	@go tool-cover -html=coverprofile=coverages/coverage.out ./internal/... ./pkg/ 2>&1 | tail -1 | \
	@echo ""
	@echo "Note: cmd/task (main package) is excluded from coverage."
	@echo "Main packages are Go conventionally excluded from unit test coverage."
	@echo "The entry points are tested via internal/cmd tests."

report:
