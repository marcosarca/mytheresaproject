# Directories to ignore (comma-separated)
IGNORE_DIRS := "docs,http,internal/config"

# Find all directories containing test files, excluding the ignored directories
TEST_DIRS := $(shell find . -type f -name '*_test.go' \
		$(foreach dir,$(IGNORE_DIRS),-not -path './$(dir)/*') \
		-exec dirname {} \; | sort -u)

# Run tests with coverage and output the coverage results
.PHONY: test
test:
	@echo "Running tests with coverage..."
	@for dir in $(TEST_DIRS); do \
		echo "Testing $$dir:"; \
		go test -coverprofile=coverage.out $$dir; \
		go tool cover -func=coverage.out | grep -E 'total:.*' | tee -a coverage_results.txt; \
		rm -f coverage.out; \
	done
	@echo "Coverage results saved to coverage_results.txt"
