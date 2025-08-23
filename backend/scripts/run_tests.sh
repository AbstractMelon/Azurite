#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COVERAGE_THRESHOLD=70
COVERAGE_DIR="coverage"
TEST_TIMEOUT="10m"

# Helper functions
log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Create coverage directory
create_coverage_dir() {
    mkdir -p ${COVERAGE_DIR}
    rm -f ${COVERAGE_DIR}/*.out ${COVERAGE_DIR}/*.html
}

# Run basic tests
run_basic_tests() {
    log "Running basic tests..."
    if go test -timeout=${TEST_TIMEOUT} -v ./...; then
        success "All tests passed!"
        return 0
    else
        error "Some tests failed!"
        return 1
    fi
}

# Run tests with race detection
run_race_tests() {
    log "Running tests with race detection..."
    if go test -timeout=${TEST_TIMEOUT} -v -race ./...; then
        success "Race tests passed!"
        return 0
    else
        error "Race condition detected!"
        return 1
    fi
}

# Run tests with coverage
run_coverage_tests() {
    log "Running tests with coverage..."

    # Run tests with coverage
    if go test -timeout=${TEST_TIMEOUT} -v -coverprofile=${COVERAGE_DIR}/coverage.out ./...; then
        success "Coverage tests completed!"
    else
        error "Coverage tests failed!"
        return 1
    fi

    # Generate HTML coverage report
    go tool cover -html=${COVERAGE_DIR}/coverage.out -o ${COVERAGE_DIR}/coverage.html
    log "Coverage report generated: ${COVERAGE_DIR}/coverage.html"

    # Show coverage summary
    echo ""
    log "Coverage Summary:"
    go tool cover -func=${COVERAGE_DIR}/coverage.out

    return 0
}

# Run detailed coverage analysis
run_detailed_coverage() {
    log "Running detailed coverage analysis..."

    echo "mode: atomic" > ${COVERAGE_DIR}/detailed_coverage.out

    # Get list of packages
    packages=$(go list ./... | grep -v /vendor/)

    for pkg in $packages; do
        log "Testing package: $pkg"
        go test -timeout=${TEST_TIMEOUT} -covermode=atomic -coverprofile=${COVERAGE_DIR}/profile.out "$pkg"

        if [ -f ${COVERAGE_DIR}/profile.out ]; then
            grep -v "mode: atomic" ${COVERAGE_DIR}/profile.out >> ${COVERAGE_DIR}/detailed_coverage.out
            rm ${COVERAGE_DIR}/profile.out
        fi
    done

    # Generate detailed HTML report
    go tool cover -html=${COVERAGE_DIR}/detailed_coverage.out -o ${COVERAGE_DIR}/detailed_coverage.html

    # Show detailed coverage summary
    echo ""
    log "Detailed Coverage Summary:"
    go tool cover -func=${COVERAGE_DIR}/detailed_coverage.out

    success "Detailed coverage analysis completed!"
}

# Check coverage threshold
check_coverage_threshold() {
    log "Checking coverage threshold (${COVERAGE_THRESHOLD}%)..."

    coverage_file=${COVERAGE_DIR}/coverage.out
    if [ ! -f "$coverage_file" ]; then
        error "Coverage file not found. Run coverage tests first."
        return 1
    fi

    # Extract total coverage percentage
    total_coverage=$(go tool cover -func="$coverage_file" | tail -1 | awk '{print $3}' | sed 's/%//')

    if [ -z "$total_coverage" ]; then
        error "Could not extract coverage percentage"
        return 1
    fi

    # Compare coverage with threshold
    if (( $(echo "$total_coverage >= $COVERAGE_THRESHOLD" | bc -l) )); then
        success "Coverage ${total_coverage}% meets threshold of ${COVERAGE_THRESHOLD}%"
        return 0
    else
        error "Coverage ${total_coverage}% is below threshold of ${COVERAGE_THRESHOLD}%"
        return 1
    fi
}

# Run benchmark tests
run_benchmarks() {
    log "Running benchmark tests..."

    if go test -bench=. -benchmem -run=^$ ./...; then
        success "Benchmark tests completed!"
        return 0
    else
        warning "Some benchmark tests may have failed"
        return 1
    fi
}

# Generate test report
generate_test_report() {
    log "Generating test report..."

    report_file="${COVERAGE_DIR}/test_report.txt"

    {
        echo "===================="
        echo "  AZURITE TEST REPORT"
        echo "===================="
        echo "Generated: $(date)"
        echo ""

        if [ -f "${COVERAGE_DIR}/coverage.out" ]; then
            echo "COVERAGE SUMMARY:"
            go tool cover -func=${COVERAGE_DIR}/coverage.out | tail -10
            echo ""
        fi

        echo "TEST PACKAGES:"
        go list ./... | grep -v /vendor/
        echo ""

        echo "GO VERSION:"
        go version
        echo ""

        echo "DEPENDENCIES:"
        go list -m all | head -20
    } > "$report_file"

    success "Test report generated: $report_file"
}

# Clean up test artifacts
cleanup() {
    log "Cleaning up test artifacts..."

    find . -name "*.test" -delete
    find . -name "test.db*" -delete
    rm -f coverage.out coverage.html

    success "Cleanup completed!"
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -b, --basic             Run basic tests only"
    echo "  -r, --race              Run tests with race detection"
    echo "  -c, --coverage          Run tests with coverage"
    echo "  -d, --detailed          Run detailed coverage analysis"
    echo "  -t, --threshold         Check coverage threshold"
    echo "  -bench, --benchmark     Run benchmark tests"
    echo "  -a, --all              Run all tests (default)"
    echo "  --report               Generate test report"
    echo "  --cleanup              Clean up test artifacts"
    echo ""
    echo "Environment Variables:"
    echo "  COVERAGE_THRESHOLD      Coverage threshold percentage (default: 70)"
    echo "  TEST_TIMEOUT           Test timeout duration (default: 10m)"
}

# Main execution
main() {
    local run_basic=false
    local run_race=false
    local run_coverage=false
    local run_detailed=false
    local check_threshold=false
    local run_bench=false
    local generate_report=false
    local run_cleanup=false
    local run_all=true

    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -b|--basic)
                run_basic=true
                run_all=false
                shift
                ;;
            -r|--race)
                run_race=true
                run_all=false
                shift
                ;;
            -c|--coverage)
                run_coverage=true
                run_all=false
                shift
                ;;
            -d|--detailed)
                run_detailed=true
                run_all=false
                shift
                ;;
            -t|--threshold)
                check_threshold=true
                run_all=false
                shift
                ;;
            -bench|--benchmark)
                run_bench=true
                run_all=false
                shift
                ;;
            -a|--all)
                run_all=true
                shift
                ;;
            --report)
                generate_report=true
                run_all=false
                shift
                ;;
            --cleanup)
                run_cleanup=true
                run_all=false
                shift
                ;;
            *)
                error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done

    # Check if we're in the right directory
    if [ ! -f "go.mod" ]; then
        error "Please run this script from the project root directory (where go.mod is located)"
        exit 1
    fi

    # Set coverage threshold from environment if provided
    if [ -n "$COVERAGE_THRESHOLD_ENV" ]; then
        COVERAGE_THRESHOLD=$COVERAGE_THRESHOLD_ENV
    fi

    # Create coverage directory
    create_coverage_dir

    local exit_code=0

    # Execute based on options
    if [ "$run_cleanup" = true ]; then
        cleanup
    elif [ "$generate_report" = true ]; then
        generate_test_report
    elif [ "$run_all" = true ]; then
        log "Running comprehensive test suite..."

        run_basic_tests || exit_code=$?

        if [ $exit_code -eq 0 ]; then
            run_race_tests || exit_code=$?
        fi

        if [ $exit_code -eq 0 ]; then
            run_coverage_tests || exit_code=$?
        fi

        if [ $exit_code -eq 0 ]; then
            check_coverage_threshold || exit_code=$?
        fi

        run_benchmarks || warning "Benchmark tests had issues"

        generate_test_report

        if [ $exit_code -eq 0 ]; then
            success "All tests completed successfully!"
        else
            error "Some tests failed. Check the output above."
        fi
    else
        # Run individual test types
        [ "$run_basic" = true ] && (run_basic_tests || exit_code=$?)
        [ "$run_race" = true ] && (run_race_tests || exit_code=$?)
        [ "$run_coverage" = true ] && (run_coverage_tests || exit_code=$?)
        [ "$run_detailed" = true ] && (run_detailed_coverage || exit_code=$?)
        [ "$check_threshold" = true ] && (check_coverage_threshold || exit_code=$?)
        [ "$run_bench" = true ] && (run_benchmarks || exit_code=$?)
    fi

    exit $exit_code
}

# Check for required tools
check_requirements() {
    local missing_tools=()

    if ! command -v go >/dev/null 2>&1; then
        missing_tools+=("go")
    fi

    if ! command -v bc >/dev/null 2>&1; then
        missing_tools+=("bc")
    fi

    if [ ${#missing_tools[@]} -gt 0 ]; then
        error "Missing required tools: ${missing_tools[*]}"
        error "Please install the missing tools and try again."
        exit 1
    fi
}

# Check requirements before running
check_requirements

# Run main function
main "$@"
