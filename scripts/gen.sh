#!/bin/bash

# Hz code generation script for scaffold
# Usage:
#   ./scripts/gen.sh hz-new <proto_file>    - Initial generation
#   ./scripts/gen.sh hz-update <proto_file>  - Incremental update

set -e

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
PROJECT_DIR=$(dirname "$SCRIPT_DIR")
IDL_DIR="$PROJECT_DIR/idl"
GEN_DIR="$PROJECT_DIR/gen/http"

check_hz() {
    if ! command -v hz &> /dev/null; then
        echo "Error: hz tool not found. Install with:"
        echo "  go install github.com/cloudwego/hertz/cmd/hz@latest"
        exit 1
    fi
}

hz_new() {
    local proto_file=$1
    if [ -z "$proto_file" ]; then
        echo "Usage: $0 hz-new <proto_file>"
        exit 1
    fi

    check_hz

    echo "Generating initial Hz code from $proto_file..."
    cd "$PROJECT_DIR"
    hz new \
        -idl "$IDL_DIR/$proto_file" \
        -module scaffold \
        -model_dir gen/http/model \
        -handler_dir gen/http/handler \
        -router_dir gen/http/router \
        --proto_path "$IDL_DIR"

    echo "Generation complete. Check $GEN_DIR for generated files."
}

hz_update() {
    local proto_file=$1
    if [ -z "$proto_file" ]; then
        echo "Usage: $0 hz-update <proto_file>"
        exit 1
    fi

    check_hz

    echo "Updating Hz code from $proto_file..."
    cd "$PROJECT_DIR"
    hz update \
        -idl "$IDL_DIR/$proto_file" \
        -module scaffold \
        -model_dir gen/http/model \
        -handler_dir gen/http/handler \
        --proto_path "$IDL_DIR"

    echo "Update complete."
}

case "$1" in
    hz-new)
        hz_new "$2"
        ;;
    hz-update)
        hz_update "$2"
        ;;
    *)
        echo "Usage: $0 {hz-new|hz-update} <proto_file>"
        echo ""
        echo "Commands:"
        echo "  hz-new    <proto>  - Generate initial Hz handler/router/model code"
        echo "  hz-update <proto>  - Incrementally update generated code"
        echo ""
        echo "Examples:"
        echo "  $0 hz-new scaffold.proto"
        echo "  $0 hz-update scaffold.proto"
        exit 1
        ;;
esac
