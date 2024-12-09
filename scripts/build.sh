#!/bin/bash

# Define the target operating systems and architectures with file extensions
TARGETS=(
  "linux/amd64:tar.gz"
  "linux/arm64:tar.gz"
  "darwin/amd64:tar.gz"
)

# Set the output directory
OUTPUT_DIR="../builds"

# Create the output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Build for each target
for target in "${TARGETS[@]}"; do
  # Split the target into GOOS and GOARCH with file extension
  goos_arch="${target%:*}"
  file_extension="${target#*:}"
  IFS='/' read -r goos goarch <<< "$goos_arch"

  # Set the environment variables
  export GOOS="$goos"
  export GOARCH="$goarch"

  # Build the package
  echo "Building for $GOOS/$GOARCH..."
  file_name="spoty-${GOOS}-${GOARCH}"
  output_file="$OUTPUT_DIR/$file_name"

  GOOS=$GOOS GOARCH=$GOARCH go build -o "${output_file}" ../cmd/web
        tar -czvf "${output_file}.tar.gz" "${output_file}"

    rm "${output_file}"
done