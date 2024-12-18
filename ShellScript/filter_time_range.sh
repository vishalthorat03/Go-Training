#!/bin/bash

# Usage: ./filter_time_range.sh <start_time> <end_time> <log_file>
start_time=$1
end_time=$2
log_file=$3

if [ -z "$start_time" ] || [ -z "$end_time" ] || [ -z "$log_file" ]; then
  echo "Usage: $0 <start_time> <end_time> <log_file>"
  exit 1
fi

awk -v start="$start_time" -v end="$end_time" \
'{
  split($1, timestamp, "T");
  if (timestamp[1] >= start && timestamp[1] <= end) print $0;
}' "$log_file"
