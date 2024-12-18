#!/bin/bash

# Usage: ./generate_report.sh <log_file>
log_file=$1
output_file="log_report.txt"

if [ -z "$log_file" ]; then
  echo "Usage: $0 <log_file>"
  exit 1
fi

echo "Log Report - $(date)" > "$output_file"
echo "======================" >> "$output_file"
echo "Total DEBUG Messages: $(grep -c '\[DEBUG\]' "$log_file")" >> "$output_file"
echo "Total INFO Messages: $(grep -c '\[INFO\]' "$log_file")" >> "$output_file"
echo "Total WARN Messages: $(grep -c '\[WARN\]' "$log_file")" >> "$output_file"
echo "" >> "$output_file"
echo "Detailed Logs:" >> "$output_file"
cat "$log_file" >> "$output_file"

echo "Report saved to $output_file"
