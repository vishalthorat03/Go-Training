#!/bin/bash

# Usage: ./extract_logs.sh <log_level> <log_file>
log_level=$1
log_file=$2

if [ -z "$log_level" ] || [ -z "$log_file" ]; then
  echo "Usage: $0 <log_level> <log_file>"
  exit 1
fi

grep "\[$log_level\]" "$log_file"
