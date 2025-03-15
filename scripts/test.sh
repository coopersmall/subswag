#!/bin/bash

path="$(pwd)/scripts"
main_dir=$(dirname "$path")

# ## Dependencies
# - https://github.com/cespare/reflex
# - https://github.com/rakyll/gotest@latest

if ! command -v reflex &> /dev/null
then
    echo "Installing reflex..."
    go get github.com/cespare/reflex
    echo "Installation complete!"
    echo "Make sure your \$GOPATH/bin is in your \$PATH."
fi

if ! command -v gotest &> /dev/null
then
    echo "Installing gotest..."
    go get github.com/rakyll/gotest
    echo "Installation complete!"
    echo "Make sure your \$GOPATH/bin is in your \$PATH."
fi

# Function to run all tests
run_all_tests() {
    echo "Running all tests..."
    gotest $main_dir/...
}

# Function to run tests for a specific file
run_test() {
    file=$1
    package=$(dirname "$file")
    base_name=$(basename "$file" .go)
    
    # If it's not a test file, run the associated test file
    test_file="${package}/${base_name}_test.go"
    echo "Change detected in $file, running all tests in ./$package..."
    gotest "./$package"
}

print_usage() {
    echo
    echo "Watching for changes in Go files..."
    echo "Press 'r' at any time to run all tests."
    echo "Press 'q' to quit."
    echo
}

# Export the functions so they're available to subshells
export -f run_all_tests
export -f run_test

# Run all tests initially
run_all_tests

# Use reflex to watch for changes in Go files
reflex -d none -r '\.go$' -- bash -c 'run_test "$0"' {} &

print_usage

# Function to handle user input
handle_input() {
    while true; do
        read -n 1 -s key
        case $key in
            r)
                echo "Running all tests..."
                run_all_tests
                print_usage
                ;;
            q)
                echo "Quitting..."
                kill $(jobs -p)
                exit 0
                ;;
        esac
    done
}

# Start handling user input
handle_input
