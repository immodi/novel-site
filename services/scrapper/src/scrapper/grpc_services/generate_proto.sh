#!/usr/bin/env bash
# generate_grpc.sh
# This script generates Python gRPC code from all .proto files in the proto directory

set -e

# Paths
PROTO_SRC_DIR="proto"
PROTO_OUT_DIR="../../"
PROTO_MAPPING="scrapper/grpc_services=proto"

echo "🔹 Generating gRPC code from all .proto files in ${PROTO_SRC_DIR}..."

# Loop over all .proto files
for proto_file in ${PROTO_SRC_DIR}/*.proto; do
  echo "⚙️  Processing ${proto_file}..."
  uv run python -m grpc_tools.protoc \
    -I${PROTO_MAPPING} \
    --python_out=${PROTO_OUT_DIR} \
    --grpc_python_out=${PROTO_OUT_DIR} \
    ${proto_file}
done

echo "✅ All gRPC code generated in ${PROTO_OUT_DIR}"
