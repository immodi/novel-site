#!/usr/bin/env bash
# generate_grpc.sh
# This script generates Python gRPC code from the scrapper.proto file

set -e

# Paths
PROTO_SRC_DIR="proto"
PROTO_OUT_DIR="../../"
PROTO_MAPPING="scrapper/grpc=proto"
PROTO_FILE="${PROTO_SRC_DIR}/scrapper.proto"

echo "🔹 Generating gRPC code from ${PROTO_FILE}..."

uv run python -m grpc_tools.protoc \
  -I${PROTO_MAPPING} \
  --python_out=${PROTO_OUT_DIR} \
  --grpc_python_out=${PROTO_OUT_DIR} \
  ${PROTO_FILE}

echo "✅ gRPC code generated in ${PROTO_OUT_DIR}"

