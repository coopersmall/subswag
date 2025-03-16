#!/bin/bash

path="$(pwd)/scripts"
main_dir=$(dirname "$path")
gh_actions_dir="$main_dir/.github/actions"

compile_gh_actions() {
    # Check if the actions directory exists
    if [ ! -d "$gh_actions_dir" ]; then
        echo "Actions directory not found: $gh_actions_dir"
        return 1
    fi

    # Loop through each subdirectory in the actions directory
    for dir in "$gh_actions_dir"/*; do
        # Skip if not a directory or if it's node_modules
        if [ ! -d "$dir" ] || [ "$(basename "$dir")" = "node_modules" ]; then
            continue
        fi

        echo "Compiling $dir"

        # Copy necessary files
        cp -r "$gh_actions_dir/node_modules" "$dir/"
        cp -r "$gh_actions_dir/package.json" "$dir/"
        cp -r "$gh_actions_dir/package-lock.json" "$dir/"

        # Build using ncc
        ncc build "$dir/index.js" -o "$dir/dist"

        # Clean up
        rm -rf "$dir/node_modules"
        rm -rf "$dir/package.json"
        rm -rf "$dir/package-lock.json"
    done
}

main() {
    compile_gh_actions
    exit $?
}

main

