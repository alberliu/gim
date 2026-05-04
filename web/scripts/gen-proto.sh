#!/usr/bin/env bash
# Generate TypeScript bindings for all *.proto files under ../pkg/proto.
set -e

ROOT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
PROTO_DIR="$ROOT_DIR/pkg/proto"
OUT_DIR="$(cd "$(dirname "$0")/.." && pwd)/src/gen/proto"

rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

PATH="$(cd "$(dirname "$0")/../node_modules/.bin" && pwd):$PATH" \
  protoc -I "$PROTO_DIR" \
  --es_out="$OUT_DIR" \
  --es_opt=target=ts,import_extension=js \
  $(find "$PROTO_DIR" -name '*.proto')

echo "generated under $OUT_DIR"
